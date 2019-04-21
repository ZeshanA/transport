package main

import (
	"database/sql"
	"log"
	"transport/lib/bus"
	"transport/lib/database"

	"github.com/lib/pq"
)

// Parses and stores data when notified that data has been received
func store(liveVehicleData *[]bus.VehicleJourney, dataIncoming chan bool) {
	db := database.OpenDBConnection()
	for {
		<-dataIncoming
		log.Printf("Vehicle entries received: %d\n", len(*liveVehicleData))
		insert(db, *liveVehicleData)
		log.Println("Finished sending vehicle entries to DB")
	}
}

// Batch inserts all vehicle entries in `vehicleActivity` into the DB
func insert(db *sql.DB, vehicleJourneys []bus.VehicleJourney) {
	// Start transaction
	transaction := database.CreateTransaction(db)
	stmt := createStatement(transaction)
	// Add all vehicle journeys to the insertion statement
	addEntriesToStatement(vehicleJourneys, stmt)
	database.CommitTransaction(stmt, transaction)
}

// Creates an SQL statement for batch insertion into the `arrivals` table
func createStatement(txn *sql.Tx) *sql.Stmt {
	table := database.VehicleJourneyTable
	// Prepare insertion statement
	stmt, err := txn.Prepare(pq.CopyIn(
		table.Name,
		table.Columns...,
	))
	if err != nil {
		log.Fatal(err)
	}
	return stmt
}

// Adds an insertion statement for each vehicle activity entry in `vehicleActivity` into `stmt`
func addEntriesToStatement(vehicleJourneys []bus.VehicleJourney, stmt *sql.Stmt) {
	for _, j := range vehicleJourneys {
		// Construct a DB row from each vehicle activity entry and insert the row into the DB
		_, err := stmt.Exec(j.Value()...)
		if err != nil {
			log.Printf("error occurred whilst executing insert statement for %v:\n%v\n", j, err)
		}
	}
}
