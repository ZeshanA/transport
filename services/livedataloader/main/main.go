package main

import (
	"log"
	"net/http"
)

func main() {
	initialiseServer()
}

func initialiseServer() {
	http.HandleFunc("/api/liveData", liveDataRequestHandler)
	log.Fatal(http.ListenAndServe(":8001", nil))
}

func liveDataRequestHandler(w http.ResponseWriter, _ *http.Request) {
	_, err := w.Write([]byte("Hello, world!"))
	if err != nil {
		log.Printf("error occurred in liveDataRequestHandler: %s\n", err)
	}
}
