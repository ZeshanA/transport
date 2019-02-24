package main

import (
	"database/sql"
	"log"
	"math"
	"strconv"
	"strings"
	"time"
	"transport/lib/progress"

	"transport/lib/database"

	"github.com/lib/pq"
)

var columnNames = []string{
	"line_ref",
	"direction_ref",
	"trip_id",
	"published_line_name",
	"operator_ref",
	"origin_ref",
	"destination_ref",
	"origin_aimed_departure_time",
	"situation_ref",
	"longitude",
	"latitude",
	"progress_rate",
	"occupancy",
	"vehicle_ref",
	"expected_arrival_time",
	"expected_departure_time",
	"distance_from_stop",
	"number_of_stops_away",
	"stop_point_ref",
	"timestamp",
}

func store(entries []ArrivalEntry) {
	// Open DB connection
	db := database.OpenDBConnection()
	defer db.Close()

	// Start transaction
	transaction, err := db.Begin()
	if err != nil {
		log.Fatal(err)
	}

	// Copy all entries into the DB (as part of the transaction)
	copyAllEntries(transaction, entries)

	// Commit transaction
	err = transaction.Commit()
	if err != nil {
		log.Fatal(err)
	}
}

func copyAllEntries(transaction *sql.Tx, entries []ArrivalEntry) {
	// Create Copy statement for all columns of the table
	tableName := string(database.VehicleJourneyTable)
	statement, err := transaction.Prepare(pq.CopyIn(tableName, columnNames...))
	if err != nil {
		log.Fatal(err)
	}

	// Execute Copy statement for each ArrivalEntry
	for i, entry := range entries {
		operatorRef := extractOperatorRef(entry.RouteID)
		progress.PrintAtIntervals(i, len(entries), "Inserting into DB:")

		// Execute statement with converted or naked fields (depending on if conversion is needed)
		_, err = statement.Exec(
			entry.RouteID, entry.DirectionID, entry.TripID, nil, operatorRef, nil, nil, nil, nil,
			entry.Longitude, entry.Latitude, phaseToProgressRate(entry.Phase), nil,
			vehicleIDToRef(operatorRef, entry.VehicleID), nil, nil, convertDistance(entry.NextStopDistance),
			-1, entry.NextStopID, entry.Timestamp.Format(time.RFC3339),
		)

		if err != nil {
			log.Fatal(err)
		}
	}

	// Close statement
	err = statement.Close()
	if err != nil {
		log.Fatal(err)
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
