package main

import (
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

type DateRange struct {
	Start time.Time
	End   time.Time
}

func main() {
	dateRange := getHostDateRange()
	stopDistances := stopdistance.Get(dbConn)
	avgStopDistances := stopdistance.GetAverage(dbConn)
	dataForDates := fetch.DateRange(dbConn, dateRange.Start, dateRange.End)
	labelledJourneys := labelDataForDates(dataForDates, stopDistances, avgStopDistances)
	database.Store(database.LabelledJourneyTable, bus.ExtractEntriesFromLabelledJourney, bus.LabelledJourneyToInterface(labelledJourneys))
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
	return hostID, hostCount, DateRange{sd, ed}
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
	for _, journeysOnDate := range dataForDates {
		partitionedJourneys := bus.PartitionJourneys(journeysOnDate)
		labelledData := labels.Create(partitionedJourneys, stopDistances, avgStopDistances)
		labelledJourneys = append(labelledJourneys, labelledData...)
	}
	return labelledJourneys
}
