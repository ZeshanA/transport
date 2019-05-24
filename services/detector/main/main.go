package main

import (
	"detector/calc"
	"detector/request"
	"fmt"
	"log"
	"time"
	"transport/lib/bustime"
	"transport/lib/database"
	"transport/lib/iohelper"
)

var loc, _ = time.LoadLocation("America/New_York")

func main() {
	// Open a DB connection and schedule it to be closed after the program returns
	db := database.OpenDBConnection()
	defer db.Close()
	// Extract the journey params from the user's request
	params := getParams()
	fmt.Println(params)
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
}

func getParams() request.JourneyParams {
	// TODO: These should come from a user request; using constants for now
	routeID, directionID, fromStop, toStop := "MTA NYCT_S78", 1, "MTA_200177", "MTA_201081"
	now := time.Now().In(loc)
	arrivalTime := now.Add(4 * time.Hour)
	return request.JourneyParams{RouteID: routeID, DirectionID: directionID, FromStop: fromStop, ToStop: toStop, ArrivalTime: arrivalTime}
}
