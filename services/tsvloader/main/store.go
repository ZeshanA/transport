package main

import (
	"database/sql"
	"log"
	"transport/lib/progress"

	"transport/lib/database"

	"github.com/lib/pq"
)

var columnNames = []string{
	"latitude",
	"longitude",
	"timestamp",
	"vehicle_id",
	"distance_along",
	"direction_id",
	"phase",
	"route_id",
	"trip_id",
	"next_stop_distance",
	"next_stop_id",
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
	tableName := string(database.ArrivalsTable)
	statement, err := transaction.Prepare(pq.CopyIn(tableName, columnNames...))
	if err != nil {
		log.Fatal(err)
	}

	// Execute Copy statement for each ArrivalEntry
	for i, entry := range entries {
		progress.PrintAtIntervals(i, len(entries), "Inserting into DB:")
		_, err = statement.Exec(
			entry.Latitude, entry.Longitude, entry.Timestamp, entry.VehicleID, entry.DistanceAlong, entry.DirectionID,
			entry.Phase, entry.RouteID, entry.TripID, entry.NextStopDistance, entry.NextStopID,
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
