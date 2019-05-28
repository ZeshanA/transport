package database

import (
	"database/sql"
	"fmt"
	"log"
	"strings"
	"time"
	"transport/lib/iohelper"
	"transport/lib/progress"

	"github.com/lib/pq"

	_ "github.com/lib/pq"
)

// Constants
const (
	databaseHost = "mtadata.postgres.database.azure.com"
	databasePort = 5432
	databaseName = "postgres"
	TimeFormat   = "2006-01-02 15:04:05"
	DateFormat   = "2006-01-02"
)

var TimeLoc, _ = time.LoadLocation("America/New_York")

// DBTable type holds name and column list for each table in the DB
type DBTable struct {
	Name    string
	Columns []string
}

// Timestamp is a wrapper around time.Time to allow for a custom
// UnmarshalJSON method
type Timestamp struct {
	time.Time
}

// Custom parsing of incoming timestamps
func (t *Timestamp) UnmarshalJSON(b []byte) error {
	noQuotes := strings.Replace(string(b), "\"", "", 2)
	parsed, err := time.Parse(time.RFC3339, noQuotes)
	if err != nil {
		log.Printf("error whilst parsing Timestamp: %v", err)
	}
	*t = Timestamp{parsed}
	return nil
}

// VehicleJourneyTable contains historical movements + live vehicle movements
var (
	VehicleJourneyTable = DBTable{
		"vehicle_journey",
		[]string{
			"line_ref", "direction_ref", "trip_id", "published_line_name", "operator_ref", "origin_ref",
			"destination_ref", "origin_aimed_departure_time", "situation_ref", "longitude", "latitude", "progress_rate",
			"occupancy", "vehicle_ref", "expected_arrival_time", "expected_departure_time", "distance_from_stop",
			"number_of_stops_away", "stop_point_ref", "timestamp",
		},
	}
)

// StopDistanceTable contains pairs of stops and the distance in metres between them
var (
	StopDistanceTable = DBTable{
		"stop_distance",
		[]string{
			"route_id",
			"from_stop_id", "to_stop_id",
			"distance",
			"direction_id",
		},
	}
)

// AverageDistanceTable contains pairs of route_id and average_distance between stops along that route
var (
	AverageDistanceTable = DBTable{
		"average_stop_distance",
		[]string{
			"route_id",
			"average_distance",
		},
	}
)

// LabelledJourneyTable contains labelled movement events
var (
	LabelledJourneyTable = DBTable{
		"labelled_journey",
		[]string{
			"line_ref",
			"direction_ref",
			"operator_ref",
			"origin_ref",
			"destination_ref",
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
			"time_to_stop",
		},
	}
)

// OpenDBConnection connects you to the MTAData DB in Azure, using
// username and password combination fetched from the environment
func OpenDBConnection() *sql.DB {
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

	log.Println("DB connection opened successfully!")

	return db
}

func getDatabaseLoginDetails() (username string, password string) {
	return iohelper.GetEnv("TRANSPORT_DB_USERNAME"), iohelper.GetEnv("TRANSPORT_DB_PASSWORD")
}

// CreateTransaction starts a DB transaction and returns a pointer to it
func CreateTransaction(db *sql.DB) *sql.Tx {
	transaction, err := db.Begin()
	if err != nil {
		log.Fatal(err)
	}
	return transaction
}

// CommitTransaction runs a blank .Exec() call to flush the given statement,
// before closing the statement and committing the transaction.
func CommitTransaction(stmt *sql.Stmt, transaction *sql.Tx) {
	_, err := stmt.Exec()
	if err != nil {
		log.Printf("error whilst flushing statement: %v\n", err)
	}
	err = stmt.Close()
	if err != nil {
		log.Printf("error whilst closing insertion statement: %v\n", err)
	}
	err = transaction.Commit()
	if err != nil {
		log.Printf("error whilst committing insertion transaction to db: %v\n", err)
	}
}

// Store sends an array of entries to the specified database table.
// columnExtractor takes a single entry and outputs a slice representing the corresponding
// database row.
func Store(table DBTable, columnExtractor func(interface{}) []interface{}, entries interface{}) {
	entriesSlice := entries.([]interface{})
	// Open DB connection
	db := OpenDBConnection()
	defer db.Close()

	// Start transaction
	transaction, err := db.Begin()
	if err != nil {
		log.Fatal(err)
	}

	// Copy all entries into the DB (as part of the transaction)
	CopyIntoDB(table, columnExtractor, transaction, entriesSlice)

	// Commit transaction
	err = transaction.Commit()
	if err != nil {
		log.Fatal(err)
	}
}

func CopyIntoDB(table DBTable, columnExtractor func(interface{}) []interface{}, transaction *sql.Tx, entries []interface{}) {
	// Create Copy statement for all columns of the table
	statement, err := transaction.Prepare(pq.CopyIn(table.Name, table.Columns...))
	if err != nil {
		log.Fatal(err)
	}

	// Execute Copy statement for each ArrivalEntry
	for i, entry := range entries {
		progress.PrintAtIntervals(i, len(entries), "Inserting into DB:")
		_, err := statement.Exec(columnExtractor(entry)...)
		if err != nil {
			log.Printf("database.CopyIntoDB: error whilst executing copy statement: %s\n", err)
		}
	}

	// Close statement
	err = statement.Close()
	if err != nil {
		log.Fatal(err)
	}
}

// Fetch all raw rows from given table
func FetchAllRows(db *sql.DB, tableName string) *sql.Rows {
	rows, err := db.Query(fmt.Sprintf(`SELECT * FROM %s`, tableName))
	if err != nil {
		log.Panicf("error whilst reading from db: %s\n", err)
	}
	return rows
}
