package fetch

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"transport/lib/bus"

	"github.com/tidwall/gjson"
)

const baseURL = "http://d.zeshan.me:8090/api/v1/vehicles"

func LiveJourneys(routeID string, directionID int) (map[string]bus.VehicleJourney, error) {
	// Create GET request
	req, err := http.NewRequest("get", baseURL, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	// Build up the required query string
	query := req.URL.Query()
	query.Add("LineRef", routeID)
	query.Add("DirectionRef", string(directionID))
	// Remove encodings of digits
	queryStr := strings.Replace(query.Encode(), "%0", "", -1)
	req.URL.RawQuery = queryStr

	fmt.Println(req.URL)

	// Execute the GET request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.Status != "200 OK" {
		return nil, fmt.Errorf("error fetching live journeys: received response with status: %s", resp.Status)
	}

	// Read in the response body
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading live journeys response: %s", err)
	}

	// Unmarshal the JSON response into a slice of VehicleJourney structs
	var journeys []bus.VehicleJourney
	err = json.Unmarshal(body, &journeys)
	if err != nil {
		return nil, fmt.Errorf("error unmarshalling live journeys response: %s", err)
	}

	// Return a map of the VehicleJourney structs keyed by VehicleRef
	return bus.VehicleJourneysByVehicleRef(journeys), nil
}

func RawJourneys() (gjson.Result, error) {
	// Create GET request
	req, err := http.NewRequest("get", baseURL, nil)
	if err != nil {
		return gjson.Result{}, err
	}
	req.Header.Set("Content-Type", "application/json")

	// Execute the GET request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return gjson.Result{}, err
	}
	defer resp.Body.Close()
	if resp.Status != "200 OK" {
		return gjson.Result{}, fmt.Errorf("error fetching live journeys: received response with status: %s", resp.Status)
	}

	// Read in the response body
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return gjson.Result{}, fmt.Errorf("error reading live journeys response: %s", err)
	}

	return gjson.ParseBytes(body), nil
}
