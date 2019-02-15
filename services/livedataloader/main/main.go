package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"transport/lib/iohelper"
)

const vehicleMonitoringURL = "http://bustime.mta.info/api/siri/vehicle-monitoring.json"

var vehicleData string

func main() {
	dataFetched := make(chan bool)
	go startFetching(vehicleMonitoringURL, iohelper.GetEnv("MTA_API_KEY"), dataFetched)
	<-dataFetched
	initialiseServer()
}

func startFetching(baseURL string, key string, dataFetched chan bool) {
	URLWithKey := fmt.Sprintf("%s?key=%s&version=2", baseURL, key)

	log.Printf("Fetching initial data from URL (%s)\n", URLWithKey)

	resp, err := http.Get(URLWithKey)
	if err != nil {
		log.Fatalf("Fetching URL (%s) failed due to:\n%s\n", err)
	}
	defer iohelper.CloseSafely(resp.Body, URLWithKey)

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("Reading response from URL (%s) failed due to:\n%s\n", err)
	}

	log.Printf("Succesfully fetched initial data from URL (%s)\n", URLWithKey)
	vehicleData = string(data)
	dataFetched <- true
}

func initialiseServer() {
	log.Printf("Starting HTTP server...")
	http.HandleFunc("/api/liveData", liveDataRequestHandler)
	http.HandleFunc("/health", func(w http.ResponseWriter, _ *http.Request) {
		_, _ = w.Write([]byte("Healthy!"))
	})
	log.Fatal(http.ListenAndServe(":8001", nil))
}

func liveDataRequestHandler(w http.ResponseWriter, _ *http.Request) {
	_, err := w.Write([]byte(vehicleData))
	if err != nil {
		log.Printf("error occurred in liveDataRequestHandler: %s\n", err)
	}
}
