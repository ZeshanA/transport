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

// Fetches initial data, telling the HTTP server it can start up, and fetches new data
// at a fixed time interval
func initialiseDataFetching(key string, dataLocation *string, dataWritten chan bool) {
	URLWithKey := fmt.Sprintf("%s?key=%s&version=2", vehicleMonitoringURL, key)
	fetchInitialData(URLWithKey, dataLocation, dataWritten)
	fetchAtInterval(URLWithKey, fetchFrequency, dataLocation, dataWritten)
}

// Fetches initial data and writes to the dataWritten channel once complete
func fetchInitialData(URL string, dataLocation *string, dataWritten chan bool) {
	fetch(URL, dataLocation)
	dataWritten <- true
	log.Printf("Succesfully fetched initial data from URL (%s)\n", URL)
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
func fetchAtInterval(URL string, timeBetweenFetches time.Duration, dataLocation *string, dataWritten chan bool) {
	ticker := time.NewTicker(timeBetweenFetches)
	go func() {
		for {
			<-ticker.C
			fetch(URL, dataLocation)
			dataWritten <- true
		}
	}()
}

// Fetches the JSON object at `URL`, reads it into memory and stores it at `dataLocation`
func fetch(URL string, dataLocation *string) {
	log.Printf("Fetching data from URL (%s)\n", URL)

	// Get response from URL
	resp, err := http.Get(URL)
	if err != nil {
		log.Printf("Fetching URL (%s) failed due to:\n%s\n", err)
		return
	}
	defer iohelper.CloseSafely(resp.Body, URL)

	// Load body of response into memory
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Reading response from URL (%s) failed due to:\n%s\n", err)
		return
	}

	// Store response body at dataLocation
	*dataLocation = string(data)

	log.Printf("Completed processing of URL (%s)\n", URL)
}
