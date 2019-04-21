package main

import (
	"log"
	"net/http"
)

// Starts the HTTP server which serves the live data
func initialiseServer() {
	port := ":8001"
	log.Printf("Starting HTTP server at http://localhost%s", port)

	// Attach request handlers
	http.HandleFunc("/api/v1/vehicles", liveDataRequestHandler)
	http.HandleFunc("/health", healthEndpoint)
	http.HandleFunc("/", healthEndpoint)

	// Start HTTP server
	log.Fatal(http.ListenAndServe(port, nil))
}

func liveDataRequestHandler(w http.ResponseWriter, req *http.Request) {
	// Construct response based on currently cached data (declared in main.go)
	// and the query params from the request
	response := createVehicleDataResponse(vehicleData, req.URL.Query())

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
