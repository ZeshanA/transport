package main

import (
	"fmt"
	"labeller/fetch"
	"labeller/labels"
	"labeller/stopdistance"
	"log"
	"os"
	"strconv"
	"time"
	"transport/lib/bus"
	"transport/lib/database"
	"transport/lib/dates"
)

var dbConn = database.OpenDBConnection()
var newYorkLoc, _ = time.LoadLocation("America/New_York")

// Boundary between days is at 4am
var dayBoundary = 4

type DateRange struct {
	Start time.Time
	End   time.Time
}

func main() {
	stopDistances := stopdistance.Get(dbConn)
	avgStopDistances := stopdistance.GetAverage(dbConn)
	mode := os.Args[1]
	switch mode {
	case "range":
		dr := getHostDateRange()
		processDateRange(dr, stopDistances, avgStopDistances)
		deleteFromDB(dr)
	case "live":
		for {
			sleepUntilProcessingTime()
			dateToProcess := time.Now().In(newYorkLoc).AddDate(0, 0, -1)
			dr := DateRange{dateToProcess, dateToProcess}
			processDateRange(dr, stopDistances, avgStopDistances)
			deleteFromDB(dr)
		}
	}
}

func processDateRange(dateRange DateRange, stopDistances map[stopdistance.Key]float64, avgStopDistances map[string]int) {
	dataForDates := fetch.DateRange(dbConn, dateRange.Start, dateRange.End)
	labelledJourneys := labelDataForDates(dataForDates, stopDistances, avgStopDistances)
	database.Store(database.LabelledJourneyTable, bus.ExtractEntriesFromLabelledJourney, bus.LabelledJourneyToInterface(labelledJourneys))
}

func sleepUntilProcessingTime() {
	// Sleep until it's 4am
	t := time.Now().In(newYorkLoc)
	endDate := t.Day()
	// If the current time is after 4am, the current day will be cut off *tomorrow* at 4am
	if t.Hour() >= dayBoundary {
		endDate += 1
	}
	// Sleep until the dayBoundary time (4am) has been reached
	processAtTime := time.Date(t.Year(), t.Month(), endDate, dayBoundary, 0, 0, 0, newYorkLoc)
	timeToSleep := processAtTime.Sub(t)
	log.Printf("Sleeping until %s\n", processAtTime.Format(database.TimeFormat))
	time.Sleep(timeToSleep)
}

func getHostDateRange() DateRange {
	hostID, hostCount, totalDateRange := extractCLIArgs()
	r := calculateHostDateRange(hostID, hostCount, totalDateRange)
	log.Printf("Labelling from %s to %s\n", r.Start, r.End)
	return r
}

func extractCLIArgs() (hostID int, hostCount int, dateRange DateRange) {
	if len(os.Args) < 5 {
		log.Fatalf("Not enough arguments provided; you must include <hostID> <hostCount> <startDate> <endDate>")
	}
	hostID, hostIDErr := strconv.Atoi(os.Args[1])
	hostCount, hostCountErr := strconv.Atoi(os.Args[2])
	if hostIDErr != nil || hostCountErr != nil {
		log.Fatalf("Failed to convert one or more arguments to integers: %v\n", os.Args)
	}
	sd, sdErr := time.Parse("2006-01-02", os.Args[3])
	ed, edErr := time.Parse("2006-01-02", os.Args[4])
	if sdErr != nil || edErr != nil {
		log.Fatalf("Failed to parse start or end date from args: %v\n", os.Args)
	}
	return hostID + 1, hostCount + 1, DateRange{sd, ed}
}

// Returns the (inclusive) date range that this host needs to label.
func calculateHostDateRange(hostID int, hostCount int, dRange DateRange) DateRange {
	daysCount := dates.DaysBetween(dRange.Start, dRange.End)
	daysPerHost := daysCount / hostCount
	offset := daysPerHost * (hostID - 1)
	startDate := dRange.Start.AddDate(0, 0, offset)
	endDate := startDate.AddDate(0, 0, daysPerHost-1)
	return DateRange{startDate, endDate}
}

func labelDataForDates(dataForDates [][]bus.VehicleJourney, stopDistances map[stopdistance.Key]float64, avgStopDistances map[string]int) []bus.LabelledJourney {
	var labelledJourneys []bus.LabelledJourney
	for i, journeysOnDate := range dataForDates {
		log.Printf("Labelling data for date %d of %d...\n", i, len(dataForDates))
		partitionedJourneys := bus.PartitionJourneys(journeysOnDate)
		labelledData := labels.Create(partitionedJourneys, stopDistances, avgStopDistances)
		labelledJourneys = append(labelledJourneys, labelledData...)
	}
	log.Println("Successfully labelled data for all dates!")
	return labelledJourneys
}

func deleteFromDB(dateRange DateRange) {
	// Start transaction
	transaction, err := dbConn.Begin()
	if err != nil {
		log.Fatal(err)
	}
	start, end := dateRange.Start, dateRange.End.AddDate(0, 0, 1)
	startStamp := dates.SetHour(start, dayBoundary, newYorkLoc).Format(database.TimeFormat)
	endStamp := dates.SetHour(end, dayBoundary, newYorkLoc).Format(database.TimeFormat)
	log.Printf("Deleting entries in DB with timestamps between %s and %s", startStamp, endStamp)
	query := fmt.Sprintf(
		"DELETE FROM %s WHERE timestamp BETWEEN '%s' and '%s'",
		database.VehicleJourneyTable.Name, startStamp, endStamp,
	)
	stat, err := transaction.Prepare(query)
	if err != nil {
		log.Fatal(err)
	}
	_, err = stat.Exec()
	if err != nil {
		log.Printf("deleteFromDB: error whilst executing delete statement: %s\n", err)
	}
	err = transaction.Commit()
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Successfully deleted entries in DB with timestamps between %s and %s", startStamp, endStamp)
}
