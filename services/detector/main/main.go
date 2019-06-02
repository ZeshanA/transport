package main

import (
	"detector/calc"
	"detector/eval"
	"detector/fetch"
	"detector/monitor"
	"detector/request"
	"log"
	"os"
	"time"
	"transport/lib/bustime"
	"transport/lib/database"
	"transport/lib/iohelper"
)

func main() {
	if len(os.Args) < 2 {
		log.Fatalf("Please pass evaluation mode (-e) or server mode (-s) as a CLI argument")
	}
	mode := os.Args[1]
	if mode == "-e" {
		eval.Evaluate()
	} else {
		server()
	}
}

func server() {
	log.Println("Server mode")
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
	log.Printf("Average time: %d\n", avgTime)
	predictedTime, err := fetch.PredictedJourneyTime(params, avgTime, stopList)
	if err != nil {
		log.Fatalf("error calculating predicted time: %s", err)
	}
	log.Printf("Predicted time: %d\n", predictedTime)
	log.Printf("Time now is: %s", time.Now().In(database.TimeLoc).Format(database.TimeFormat))
	log.Printf("Arrival time is: %s", params.ArrivalTime.In(database.TimeLoc).Format(database.TimeFormat))
	complete := make(chan bool)
	monitor.LiveBuses(avgTime, predictedTime, params, stopList, db, complete)
	<-complete
}
