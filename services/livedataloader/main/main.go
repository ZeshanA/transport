package main

import (
	"transport/lib/bus"
	"transport/lib/iohelper"
)

// Currently cached data from MTA
var vehicleData []bus.VehicleJourney

func main() {
	// Create a channel which is written to when new data is finished being written
	dataIncoming := make(chan bool)
	// When new data arrives, store it in the historical DB
	go store(&vehicleData, dataIncoming)
	// Set up data polling
	initialiseDataFetching(iohelper.GetEnv("MTA_API_KEY"), &vehicleData, dataIncoming)
	// Start HTTP server
	initialiseServer()
}
