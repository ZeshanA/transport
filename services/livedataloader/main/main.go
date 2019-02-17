package main

import (
	"log"
	"net/http"
	"transport/lib/iohelper"
)

// Currently cached data from MTA
var vehicleData string

func main() {
	// Fetch initial data and set up polling
	initialiseDataFetching(iohelper.GetEnv("MTA_API_KEY"), &vehicleData)
	// Start serving!
	initialiseServer()
}

// Starts the HTTP server which serves the live data
func initialiseServer() {
	log.Printf("Starting HTTP server...")

	// Attach request handlers
	http.HandleFunc("/api/liveData", liveDataRequestHandler)
	http.HandleFunc("/health", healthEndpoint)

	// Start HTTP server
	log.Fatal(http.ListenAndServe(":8001", nil))
}

func liveDataRequestHandler(w http.ResponseWriter, req *http.Request) {
	// Construct response based on currently cached data and the query params from the request
	response := *createVehicleDataResponse(&vehicleData, req.URL.Query())

	// Write response
	_, err := w.Write([]byte(response))
	if err != nil {
		log.Printf("error occurred in liveDataRequestHandler: %s\n", err)
	}
}

func healthEndpoint(w http.ResponseWriter, _ *http.Request) {
	_, err := w.Write([]byte("Healthy!"))
	if err != nil {
		log.Printf("error occurred in healthEndpoint: %s\n", err)
	}
}
