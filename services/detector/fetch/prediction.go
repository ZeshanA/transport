package fetch

import (
	"bytes"
	"detector/request"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"
	"transport/lib/bus"
	"transport/lib/bustime"

	"github.com/tidwall/gjson"

	"gopkg.in/guregu/null.v3"
)

const predictionBaseURL = "http://127.0.0.1:5000"
const stopToStopURL = predictionBaseURL + "/predictStopToStop"
const singleMovementURL = predictionBaseURL + "/predictFromMovement"

func SingleMovementPrediction(journey bus.VehicleJourney) (int, error) {
	jsonStr, err := json.Marshal(journey)
	fmt.Println(string(jsonStr))
	if err != nil {
		log.Printf("error marshalling journey into JSON: %s", err)
	}
	req, err := http.NewRequest("POST", singleMovementURL, bytes.NewBuffer(jsonStr))
	if err != nil {
		return 0, err
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	if resp.Status != "200 OK" {
		return 0, fmt.Errorf("error fetching predicted arrival time: received response with status: %s", resp.Status)
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return 0, fmt.Errorf("error reading predicted arrival time response: %s", err)
	}
	timeStr := gjson.GetBytes(body, "prediction").String()
	time, err := strconv.Atoi(timeStr)
	if err != nil {
		return 0, fmt.Errorf("error reading predicted arrival time response: %s", err)
	}
	return time, nil
}

func PredictedJourneyTime(params request.JourneyParams, avgTime int, stopList []bustime.BusStop) (int, error) {
	var jsonStr = createJSONRequest(params, avgTime, stopList)
	req, err := http.NewRequest("POST", stopToStopURL, bytes.NewBuffer(jsonStr))
	if err != nil {
		return 0, err
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	if resp.Status != "200 OK" {
		return 0, fmt.Errorf("error fetching predicted journey time: received response with status: %s", resp.Status)
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return 0, fmt.Errorf("error reading predicted journey time response: %s", err)
	}
	timeStr := gjson.GetBytes(body, "prediction").String()
	time, err := strconv.Atoi(timeStr)
	if err != nil {
		return 0, fmt.Errorf("error reading predicted journey time response: %s", err)
	}
	return time, nil
}

func createJSONRequest(params request.JourneyParams, avgTime int, stopList []bustime.BusStop) []byte {
	template := `{
		"averageJourneyTime": %d,
		"sampleMovement": %s,
		"journey": %s,
		"stopList": %s
	}
	`
	paramsJSON, _ := params.MarshalJSON()
	stopListJSON, _ := json.Marshal(stopList)
	str := fmt.Sprintf(template, avgTime, createSampleMovement(params, stopList), paramsJSON, stopListJSON)
	return []byte(str)
}

func createSampleMovement(params request.JourneyParams, stopList []bustime.BusStop) string {
	operator := strings.Split(params.RouteID, "_")[0]
	lj := bus.VehicleJourney{
		LineRef:        null.StringFrom(params.RouteID),
		DirectionRef:   null.IntFrom(int64(params.DirectionID)),
		OperatorRef:    null.StringFrom(operator),
		OriginRef:      null.StringFrom(stopList[0].ID),
		DestinationRef: null.StringFrom(stopList[len(stopList)-1].ID),
		ProgressRate:   null.StringFrom("normalProgress"),
		VehicleRef:     null.StringFrom("MTA NYCT_7339"),
	}
	jsonStr, err := json.Marshal(lj)
	if err != nil {
		log.Printf("error marshalling sample movement: %s", err)
		return ""
	}
	return string(jsonStr)
}
