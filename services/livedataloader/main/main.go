package main

import (
	"transport/lib/iohelper"
)

// Currently cached data from MTA
var vehicleData string

func main() {
	// Create a channel which is written to when new data is finished being written
	dataIncoming := make(chan bool)
	// When new data arrives, store it in the historical DB
	go store(dataIncoming)
	// Set up data polling
	initialiseDataFetching(iohelper.GetEnv("MTA_API_KEY"), &vehicleData, dataIncoming)
	// Start HTTP server
	initialiseServer()
}
