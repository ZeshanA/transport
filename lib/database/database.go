package database

import (
	"database/sql"
	"fmt"
	"log"
	"transport/lib/iohelper"

	_ "github.com/lib/pq"
)

// Constants
const (
	databaseHost = "mtadata.postgres.database.azure.com"
	databasePort = 5432
	databaseName = "postgres"
)

// DBTable type holds strings for each table in the DB
type DBTable string

const (
	// ArrivalsTable contains historical movements + live vehicle movements
	ArrivalsTable DBTable = "arrivals"
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
