package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"
	"transport/lib/iohelper"
)

// Constants
const (
	vehicleMonitoringURL = "http://bustime.mta.info/api/siri/vehicle-monitoring.json"
	fetchFrequency       = 35 * time.Second
)

var vehicleData string

func main() {
	initialDataFetched := make(chan bool)
	go initialiseDataFetching(vehicleMonitoringURL, iohelper.GetEnv("MTA_API_KEY"), initialDataFetched)
	<-initialDataFetched
	initialiseServer()
}

// Fetches initial data, telling the HTTP server it can start up, and fetches new data
// at a fixed time interval
func initialiseDataFetching(baseURL string, key string, initialDataFetched chan bool) {
	URLWithKey := fmt.Sprintf("%s?key=%s&version=2", baseURL, key)
	fetchInitialData(URLWithKey, initialDataFetched)
	fetchAtInterval(URLWithKey, fetchFrequency)
}

// Fetches initial data and tells the HTTP server it can start up (via the `initialDataFetched` channel)
func fetchInitialData(URL string, initialDataFetched chan bool) {
	fetch(URL)
	log.Printf("Succesfully fetched initial data from URL (%s)\n", URL)
	initialDataFetched <- true
}

/*
This function fetches the data at URL at intervals of `timeBetweenFetches`.

Implementation:
Make a time.NewTicker, which returns a channel that will be written to every X seconds.
   In a new go routine, the infinite for loop:
	- blocks and waits for the ticker channel to be written to
	- fetches the data
	- returns to the start of the loop and blocks on the channel again
*/
func fetchAtInterval(URL string, timeBetweenFetches time.Duration) {
	ticker := time.NewTicker(timeBetweenFetches)
	go func() {
		for {
			<-ticker.C
			fetch(URL)
		}
	}()
}

// Fetches the JSON object at `URL`, reads it into memory and stores it at `vehicleData`
func fetch(URL string) {
	log.Printf("Fetching data from URL (%s)\n", URL)

	resp, err := http.Get(URL)
	if err != nil {
		log.Printf("Fetching URL (%s) failed due to:\n%s\n", err)
		return
	}
	defer iohelper.CloseSafely(resp.Body, URL)

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Reading response from URL (%s) failed due to:\n%s\n", err)
		return
	}

	vehicleData = string(data)

	log.Printf("Completed processing of URL (%s)\n", URL)
}

// Starts the HTTP server which serves the live data
func initialiseServer() {
	log.Printf("Starting HTTP server...")
	http.HandleFunc("/api/liveData", liveDataRequestHandler)
	http.HandleFunc("/health", func(w http.ResponseWriter, _ *http.Request) {
		_, _ = w.Write([]byte("Healthy!"))
	})
	log.Fatal(http.ListenAndServe(":8001", nil))
}

func liveDataRequestHandler(w http.ResponseWriter, req *http.Request) {
	// Construct response based on currently cached data and the query params in request
	response := createVehicleDataResponse(vehicleData, req.URL.Query())

	// Write response
	_, err := w.Write([]byte(response))
	if err != nil {
		log.Printf("error occurred in liveDataRequestHandler: %s\n", err)
	}
}
