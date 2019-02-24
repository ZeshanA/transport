package main

import (
	"errors"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"
	"transport/lib/database"

	"github.com/tidwall/gjson"
)

const timeFormat = "2006-01-02 15:04:05"

func store(liveVehicleData *string, dataIncoming chan bool) {
	database.OpenDBConnection()
	for {
		<-dataIncoming
		parseAndStore(liveVehicleData)
		log.Println("DONE SENDING ALL TO DB")
	}
}

func parseAndStore(liveVehicleData *string) {
	// Extract vehicle activity from JSON string
	vehicleActivity := gjson.Get(*liveVehicleData, VehicleActivityPath)

	// Construct a DB row from each vehicle activity entry and insert the row into the DB
	vehicleActivity.ForEach(func(_, activityEntry gjson.Result) bool {
		fieldValues, err := getFieldValues(&activityEntry)
		if err != nil {
			return true
		}
		fmt.Printf("%v\n", fieldValues)
		// Insert fields into DB
		return true
	})
}

func getFieldValues(activityEntry *gjson.Result) (*[]string, error) {
	// Some fields are nested directly under journey, others are nested more deeply
	journey := activityEntry.Get("MonitoredVehicleJourney")
	location := journey.Get("VehicleLocation")
	call := journey.Get("MonitoredCall")

	fieldCount := 11
	fields := make([]string, fieldCount)

	// Store fields in slice
	fields[0], fields[1] = location.Get("Latitude").String(), location.Get("Longitude").String()
	timestamp, err := parseTime(activityEntry.Get("RecordedAtTime").String())
	if err != nil {
		return nil, err
	}
	fields[2] = timestamp
	vehicleID, err := intFromID(journey.Get("VehicleRef").String())
	if err != nil {
		return nil, err
	}
	fields[3] = strconv.Itoa(vehicleID)
	fields[4] = "0"
	fields[5] = journey.Get("DirectionRef").String()
	phase, err := parseProgressRate(journey.Get("ProgressRate").String())
	if err != nil {
		return nil, err
	}
	fields[6] = phase
	fields[7] = journey.Get("LineRef").String()
	fields[8] = journey.Get("BlockRef").String()
	fields[9] = call.Get("DistanceFromStop").String()
	fields[10] = call.Get("StopPointRef").String()
	return &fields, nil
}

func intFromID(stringID string) (numericID int, parseError error) {
	split := strings.Split(stringID, "_")
	numericValue, err := strconv.Atoi(split[1])
	if err != nil {
		return 0, err
	}
	return numericValue, nil
}

/*
Indicator of whether the bus is:
	- making progress (normalProgress) (i.e. moving, generally),
	- not moving (with value noProgress),
	- laying over before beginning a trip (value layover),
	- or serving a trip prior to one which will arrive (prevTrip).
*/
func parseProgressRate(progressRate string) (string, error) {
	switch progressRate {
	case "normalProgress", "noProgress":
		return "IN_PROGRESS", nil
	case "layover":
		return "LAYOVER_DURING", nil
	case "prevTrip":
		return "PREV_TRIP", nil
	default:
		return "", errors.New(fmt.Sprintf("invalid progress rate: %s\n", progressRate))
	}
}

func parseTime(timeString string) (string, error) {
	parsed, err := time.Parse(time.RFC3339, timeString)
	if err != nil {
		return "", errors.New(fmt.Sprintf("invalid timestamp: %s\n", timeString))
	}
	return parsed.Format(timeFormat), nil
}
