package main

import (
	"fmt"
	"log"
	"transport/lib/iohelper"
	"transport/lib/progress"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

const (
	databaseHost = "mtadata.postgres.database.azure.com"
	databasePort = 5432
	databaseName = "postgres"
)

func getDatabaseLoginDetails() (username string, password string) {
	return iohelper.GetEnv("TRANSPORT_DB_USERNAME"), iohelper.GetEnv("TRANSPORT_DB_PASSWORD")
}

// Stores each struct in the entries array into the Postgres DB
func store(entries []ArrivalEntry) {
	// Create DB connection string using login details from environment
	username, password := getDatabaseLoginDetails()
	connectionDetails := fmt.Sprintf(
		"host=%s port=%d user=%s dbname=%s password=%s",
		databaseHost, databasePort, username, databaseName, password,
	)

	// Open DB connection
	db, err := gorm.Open("postgres", connectionDetails)
	if err != nil {
		log.Fatalf("Failed to open DB connection due to: %v\n", err)
	}
	defer db.Close()

	// Put each struct from unmarshalMTADataFile array into Postgres
	for i, entry := range entries {
		progress.PrintAtIntervals(i, len(entries), "Inserting into DB:")
		db.NewRecord(entry)
		db.Create(&entry)
	}
}
