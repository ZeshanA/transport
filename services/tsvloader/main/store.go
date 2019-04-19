package main

import (
	"log"
	"math"
	"strconv"
	"strings"
	"time"
)

// Takes an ArrivalEntry struct and outputs a slice representing the database
// columns for that entry
func extractColsFromArrivalEntry(arrivalEntry interface{}) []interface{} {
	entry, ok := arrivalEntry.(ArrivalEntry)
	if !ok {
		log.Panicf("processArrivalEntry: entry passed in is not an ArrivalEntry struct")
	}
	operatorRef := extractOperatorRef(entry.RouteID)
	return []interface{}{
		entry.RouteID, entry.DirectionID, entry.TripID, nil, operatorRef, nil, nil, nil, nil,
		entry.Longitude, entry.Latitude, phaseToProgressRate(entry.Phase), nil,
		vehicleIDToRef(operatorRef, entry.VehicleID), nil, nil, convertDistance(entry.NextStopDistance),
		-1, entry.NextStopID, entry.Timestamp.Format(time.RFC3339),
	}
}

// Extracts the operator reference from a routeID
func extractOperatorRef(routeID string) (operatorRef string) {
	split := strings.Split(routeID, "_")
	if len(split) == 0 {
		return ""
	}
	return split[0]
}

// Converts an archaic "phase" string to a modern "progressRate" string
func phaseToProgressRate(phase string) (progressRate string) {
	switch phase {
	case "IN_PROGRESS":
		return "normalProgress"
	case "LAYOVER_DURING", "LAYOVER_BEFORE":
		return "layover"
	default:
		return "normalProgress"
	}
}

// Produces a vehicle reference (e.g. "MTA NYCT_123"), given an operatorRef and a vehicleID
func vehicleIDToRef(operatorRef string, vehicleID int32) (vehicleRef string) {
	return operatorRef + "_" + strconv.Itoa(int(vehicleID))
}

// Converts float distances into integers by rounding
func convertDistance(distance float64) (converted int) {
	return int(math.Round(distance))
}
