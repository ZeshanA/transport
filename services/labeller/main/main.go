package main

import (
	"fmt"
	"labeller/fetch"
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
	fmt.Println(stopDistances)
	dataForDates := fetch.DateRange(dbConn, dateRange[0], dateRange[1])
	partitionedDataForDates := make([]map[bus.DirectedRoute][]bus.VehicleJourney, len(dataForDates))
	for i, journeysOnDate := range dataForDates {
		partitionedDataForDates[i] = bus.PartitionJourneys(journeysOnDate)
	}
	fmt.Println(partitionedDataForDates)
}
