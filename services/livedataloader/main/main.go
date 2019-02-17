package main

import (
	"transport/lib/iohelper"
)

// Currently cached data from MTA
var vehicleData string

func main() {
	// Set up data polling and start HTTP Server
	initialiseDataFetching(iohelper.GetEnv("MTA_API_KEY"), &vehicleData)
	initialiseServer()
}
