package main

import (
	"database/sql"
	"fmt"
	"log"
	"transport/lib/iohelper"
	"transport/lib/progress"

	"github.com/lib/pq"
)

const (
	databaseHost = "mtadata.postgres.database.azure.com"
	databasePort = 5432
	databaseName = "postgres"
	tableName    = "arrival2"
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

func getDatabaseLoginDetails() (username string, password string) {
	return iohelper.GetEnv("TRANSPORT_DB_USERNAME"), iohelper.GetEnv("TRANSPORT_DB_PASSWORD")
}

func store(entries []ArrivalEntry) {
	// Open DB connection
	db := openDBConnection()
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

func openDBConnection() *sql.DB {
	// Create DB connection string
	username, password := getDatabaseLoginDetails()
	connectionDetails := fmt.Sprintf(
		"host=%s port=%d user=%s dbname=%s password=%s",
		databaseHost, databasePort, username, databaseName, password,
	)

	// Connect to DB
	db, err := sql.Open("postgres", connectionDetails)
	if err != nil {
		log.Fatalf("Failed to open DB connection due to: %v\n", err)
	}

	// Set connection params
	db.SetMaxIdleConns(10)
	db.SetMaxOpenConns(10)
	db.SetConnMaxLifetime(0)

	fmt.Println("DB connection opened succesfully!")

	return db
}

func copyAllEntries(transaction *sql.Tx, entries []ArrivalEntry) {
	// Create Copy statement for all columns of the table
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
