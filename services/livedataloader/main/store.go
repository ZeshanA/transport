package main

import (
	"errors"
	"fmt"
	"log"
	"strconv"
	"strings"
	"transport/lib/database"

	"github.com/tidwall/gjson"
)

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
		_, err := getFieldValues(&activityEntry)
		if err != nil {
			return true
		}
		// Insert fields into DB
		return true
	})
}

func getFieldValues(activityEntry *gjson.Result) (*[]string, error) {
	// Some fields are nested directly under journey, others are nested more deeply
	journey := activityEntry.Get("MonitoredVehicleJourney")
	location := journey.Get("VehicleLocation")
	call := journey.Get("MonitoredCall")

	// Extract field values
	latitude, longitude := location.Get("Longitude"), location.Get("Longitude")
	timestamp := activityEntry.Get("RecordedAtTime")

	vehicleID, err := intFromID(journey.Get("VehicleRef").String())
	if err != nil {
		return nil, err
	}

	directionID := journey.Get("DirectionRef")

	phase, err := parseProgressRate(journey.Get("ProgressRate").String())
	if err != nil {
		return nil, err
	}

	routeID := journey.Get("LineRef")
	tripID := journey.Get("BlockRef")
	nextStopDistance := call.Get("DistanceFromStop")
	nextStopID := call.Get("StopPointRef")

	log.Printf("Lat: %v, Lon: %v, Timestamp: %v, VehicleID: %v, directionID: %v, next: %v, next: %v, next: %v, next: %v, next: %v\n\n\n", latitude, longitude, timestamp, vehicleID, directionID, phase, routeID, tripID, nextStopDistance, nextStopID)
	return nil, nil
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
