package main

import (
	"detector/calc"
	"detector/fetch"
	"detector/request"
	"fmt"
	"log"
	"transport/lib/bustime"
	"transport/lib/database"
	"transport/lib/iohelper"
)

func main() {
	// Open a DB connection and schedule it to be closed after the program returns
	db := database.OpenDBConnection()
	defer db.Close()
	// Extract the journey params from the user's request
	params := request.GetParams()
	// Create a bustime client to fetch list of stops
	bt := bustime.NewClient(iohelper.GetEnv("MTA_API_KEY"))
	// Fetch the list of stops for the requested route and direction
	stopList := bt.GetStops(params.RouteID)[params.RouteID][params.DirectionID]
	// Get average time to travel between stops
	avgTime, err := calc.AvgTimeBetweenStops(stopList, params, db)
	if err != nil {
		log.Fatalf("error calculating average time between stops: %s", err)
	}
	fmt.Printf("Average time: %d\n", avgTime)
	predictedTime, err := fetch.PredictedJourneyTime(params, avgTime, stopList)
	if err != nil {
		log.Fatalf("error calculating predicted time: %s", err)
	}
	fmt.Printf("Predicted time: %d\n", predictedTime)
}
