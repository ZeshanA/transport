package main

import (
	"log"
	"net/http"
)

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
	// Construct response based on currently cached data (declared in main.go)
	// and the query params from the request
	// TODO: Implement querying for the new internal VehicleJourney format
	response := vehicleData // *createVehicleDataResponse(&vehicleData, req.URL.Query())

	log.Printf("Response created succesfully, writing to output...")

	// Write response
	_, err := w.Write([]byte(response))
	if err != nil {
		log.Printf("error occurred whilst writing response in liveDataRequestHandler: %s\n", err)
	}

	log.Printf("Response completed succesfully!")
}

func healthEndpoint(w http.ResponseWriter, _ *http.Request) {
	_, err := w.Write([]byte("Healthy!"))
	if err != nil {
		log.Printf("error occurred in healthEndpoint: %s\n", err)
	}
}
