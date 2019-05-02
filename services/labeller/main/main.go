package main

import (
	"labeller/fetch"
	"labeller/labels"
	"labeller/stopdistance"
	"time"
	"transport/lib/bus"
	"transport/lib/database"
)

var dbConn = database.OpenDBConnection()
var dateRange = []time.Time{
	time.Date(2019, 04, 21, 0, 0, 0, 0, time.UTC),
	time.Date(2019, 04, 22, 0, 0, 0, 0, time.UTC),
}

func main() {
	stopDistances := stopdistance.Get(dbConn)
	avgStopDistances := stopdistance.GetAverage(dbConn)
	dataForDates := fetch.DateRange(dbConn, dateRange[0], dateRange[1])
	var labelledJourneys []bus.LabelledJourney
	for _, journeysOnDate := range dataForDates {
		partitionedJourneys := bus.PartitionJourneys(journeysOnDate)
		labelledData := labels.Create(partitionedJourneys, stopDistances, avgStopDistances)
		labelledJourneys = append(labelledJourneys, labelledData...)
	}
	//for i, j := range labelledJourneys {
	//	fmt.Printf("Journey %d: %v\n", i, j)
	//}
}
