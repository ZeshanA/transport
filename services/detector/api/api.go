package api

import (
	"database/sql"
	"detector/calc"
	"detector/fetch"
	"detector/monitor"
	"detector/request"
	"detector/response"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"time"
	"transport/lib/bustime"
	"transport/lib/database"
	"transport/lib/iohelper"
	"transport/lib/network"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

var bt = bustime.NewClient(iohelper.GetEnv("MTA_API_KEY"))
var stopInfo []byte
var parsedStopInfo map[string]map[int][]bustime.BusStop
var db *sql.DB

func Start() {
	r := mux.NewRouter()
	fetchStopDetails()
	r.HandleFunc("/getStops", fetchStops)
	r.HandleFunc("/subscribe", subscribe)
	// Open a DB connection and schedule it to be closed after the program returns
	db = database.OpenDBConnection()
	defer db.Close()
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
	parsedStopInfo = stopDetails
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

func subscribe(w http.ResponseWriter, r *http.Request) {
	params := extractParams(w, r)
	go notifyAtOptimalDeparture(params)
	msg := map[string]string{
		"status": "ok",
	}
	err := json.NewEncoder(w).Encode(msg)
	if err != nil {
		log.Printf("error writing JSON response: %v", err)
	}
}

func notifyAtOptimalDeparture(params request.JourneyParams) {
	response.SendNotification(params, response.Notification{
		VehicleID: "LOL SET OFF WHENEVER YOU WANT MATE",
	})
	// Get the list of stops for the requested route and direction
	stopList := parsedStopInfo[params.RouteID][params.DirectionID]
	// Get average time to travel between stops
	avgTime, err := calc.AvgTimeBetweenStops(stopList, params, db)
	if err != nil {
		log.Fatalf("error calculating average time between stops: %s", err)
	}
	log.Printf("Average time: %d\n", avgTime)
	predictedTime, err := fetch.PredictedJourneyTime(params, avgTime, stopList)
	if err != nil {
		log.Fatalf("error calculating predicted time: %s", err)
	}
	log.Printf("Predicted time: %d\n", predictedTime)
	log.Printf("Time now is: %s", time.Now().In(database.TimeLoc).Format(database.TimeFormat))
	log.Printf("Arrival time is: %s", params.ArrivalTime.In(database.TimeLoc).Format(database.TimeFormat))
	complete := make(chan response.Notification)
	monitor.LiveBuses(avgTime, predictedTime, params, stopList, db, complete)
	notification := <-complete
	response.SendNotification(params, notification)
}

func extractParams(w http.ResponseWriter, r *http.Request) request.JourneyParams {
	jsonBytes, err := ioutil.ReadAll(r.Body)
	if err != nil {
		network.WriteError("error reading JSON body from response: %v", err, w)
	}
	var params request.JourneyParams
	err = json.Unmarshal(jsonBytes, &params)
	if err != nil {
		network.WriteError("error unmarshalling JSON: %v", err, w)
	}
	return params
}
