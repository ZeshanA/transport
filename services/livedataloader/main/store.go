package main

import (
	"database/sql"
	"database/sql/driver"
	"log"
	"transport/lib/database"

	"github.com/lib/pq"
	"github.com/tidwall/gjson"
)

// Parses and stores data when notified that data has been received
func store(liveVehicleData *string, dataIncoming chan bool) {
	db := database.OpenDBConnection()
	for {
		<-dataIncoming
		parseAndStore(liveVehicleData, db)
		log.Println("Finished sending vehicle entries to DB")
	}
}

// Parses string data into gjson Result, extracts relevant fields and inserts into DB
func parseAndStore(liveVehicleData *string, db *sql.DB) {
	// Extract vehicle activity from JSON string
	vehicleJourneys := gjson.Parse(*liveVehicleData)
	printEntryCount(&vehicleJourneys)
	insert(db, &vehicleJourneys)
}

// Prints the number of individual entries stored within a vehicleActivity item
func printEntryCount(vehicleJourneys *gjson.Result) {
	arr := (*vehicleJourneys).Array()
	log.Printf("Vehicle entries received: %d\n", len(arr))
}

// Batch inserts all vehicle entries in `vehicleActivity` into the DB
func insert(db *sql.DB, vehicleJourneys *gjson.Result) {
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
func addEntriesToStatement(vehicleJourneys *gjson.Result, stmt *sql.Stmt) {
	// Construct a DB row from each vehicle activity entry and insert the row into the DB
	(*vehicleJourneys).ForEach(func(_, activityEntry gjson.Result) bool {
		// Get field values
		fieldValues := getFieldValues(&activityEntry)
		// Add field values of current vehicle entry to the SQL statement
		_, err := stmt.Exec(*fieldValues...)
		if err != nil {
			log.Printf("error occurred whilst executing insert statement for %v:\n%v\n", fieldValues, err)
			return true
		}
		return true
	})
}

// Returns a slice representing the single, given `activityEntry` as a DB row
// We return an []interface to allow
func getFieldValues(vehicleJourney *gjson.Result) *[]interface{} {
	i, fieldCount := 0, 20
	fields := make([]interface{}, fieldCount)

	// Insert each field in the journey into the `fields` array
	vehicleJourney.ForEach(func(key, field gjson.Result) bool {
		// SituationRefs are a text[] in the DB, so we need to convert them to pq arrays
		// before insertion
		if key.String() == "SituationRef" {
			fields[i] = pqArrayFromJSONArray(field)
		} else {
			fields[i] = field.String()
		}
		i++
		return true
	})

	return &fields
}

// Converts a JSON array (represented as a gjson.Result struct)
// into a pq array ready for DB insertion
func pqArrayFromJSONArray(jsonArray gjson.Result) interface {
	driver.Valuer
	sql.Scanner
} {
	return pq.Array(toStringSlice(jsonArray.Array()))
}

// Converts an array of gjson.Result structs into a slice of their string values
func toStringSlice(arr []gjson.Result) []string {
	result := make([]string, len(arr))
	for i, field := range arr {
		result[i] = field.String()
	}
	return result
}
