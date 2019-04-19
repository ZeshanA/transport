package main

import (
	"database/sql"
	"fmt"
	"labeller/stopdistance"
	"time"
	"transport/lib/database"
	"transport/lib/dates"
)

var dbConn = database.OpenDBConnection()

func main() {
	stopDistances := stopdistance.Get(dbConn)
	fmt.Println(stopDistances)
}

func ProcessDateRange(db *sql.DB, startDate time.Time, lastDate time.Time) {
	endDate := lastDate.AddDate(0, 0, 1)
	for d := startDate; !dates.Equal(d, endDate); d = d.AddDate(0, 0, 1) {
		fmt.Println(d)
	}
}
