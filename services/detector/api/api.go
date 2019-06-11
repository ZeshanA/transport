package api

import (
	"encoding/json"
	"log"
	"net/http"
	"transport/lib/bustime"
	"transport/lib/iohelper"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

var bt = bustime.NewClient(iohelper.GetEnv("MTA_API_KEY"))
var stopInfo []byte

func Start() {
	r := mux.NewRouter()
	fetchStopDetails()
	r.HandleFunc("/getStops", fetchStops)
	err := http.ListenAndServe(":7891", handlers.CORS()(r))
	if err != nil {
		log.Fatalf("http server crashed with error: %v", err)
	}
}

func fetchStopDetails() {
	agencies := bt.GetAgencies()
	log.Printf("%d agencies fetched\n", len(agencies))
	routes := bt.GetRoutes(agencies...)
	log.Printf("%d routes fetched\n", len(agencies))
	stopDetails := bt.GetStops(routes...)
	jsonStopDetails, err := json.Marshal(stopDetails)
	if err != nil {
		log.Fatalf("failed to convert stop details into JSON due to error: %v", err)
	}
	stopInfo = jsonStopDetails
}

func fetchStops(w http.ResponseWriter, r *http.Request) {
	// Enable CORS
	w.Header().Set("Access-Control-Allow-Origin", "*")
	if stopInfo == nil {
		w.Write([]byte("Stops not yet fetched"))
		return
	}
	w.Write(stopInfo)
}
