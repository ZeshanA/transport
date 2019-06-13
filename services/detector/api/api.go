package api

import (
	"database/sql"
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

	"github.com/rs/cors"

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
	r.HandleFunc("/subscribe", subscribe).Methods("POST")
	// Open a DB connection and schedule it to be closed after the program returns
	db = database.OpenDBConnection()
	defer db.Close()
	corsOptions := cors.New(cors.Options{
		AllowedOrigins: []string{"*"},
	})
	handler := corsOptions.Handler(r)
	port := ":7891"
	log.Printf("HTTP server started on %s...\n", port)
	log.Fatal(http.ListenAndServe(port, handler))
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
	if stopInfo == nil {
		w.Write([]byte("Stops not yet fetched"))
		return
	}
	w.Write(stopInfo)
}

func subscribe(w http.ResponseWriter, r *http.Request) {
	params := extractParams(w, r)
	log.Printf("Received subscription request: %v", params)
	go notifyAtOptimalDeparture(params)
	msg := map[string]string{
		"status": "ok",
	}
	err := json.NewEncoder(w).Encode(msg)
	if err != nil {
		log.Printf("error writing JSON response: %v", err)
	}
	log.Printf("Succesfully registered subscription")
}

func notifyAtOptimalDeparture(params request.JourneyParams) {
	// Mock notification code
	//time.Sleep(10 * time.Second)
	//t := time.Now().In(database.TimeLoc)
	//response.SendNotification(params, response.Notification{
	//	VehicleID:            "MTA_12389",
	//	OptimalDepartureTime: database.Timestamp{Time: t},
	//	PredictedArrivalTime: database.Timestamp{Time: t.Add(20 * time.Minute)},
	//})
	// Get the list of stops for the requested route and direction
	log.Println("Extracting list of stops from cache...")
	stopList := parsedStopInfo[params.RouteID][params.DirectionID]
	// Get average time to travel between stops
	// avgTime, err := calc.AvgTimeBetweenStops(stopList, params, db)
	avgTime := 1039
	var err error
	if err != nil {
		log.Fatalf("error calculating average time between stops: %s", err)
	}
	log.Printf("Average time: %d\n", avgTime)
	log.Println("Fetching predicted journey time from predictor service...")
	predictedTime, err := fetch.PredictedJourneyTime(params, avgTime, stopList)
	if err != nil {
		log.Fatalf("error calculating predicted time: %s", err)
	}
	log.Printf("Predicted time: %d\n", predictedTime)
	log.Printf("Time now is: %s", time.Now().In(database.TimeLoc).Format(database.TimeFormat))
	log.Printf("Arrival time is: %s", params.ArrivalTime.Format(database.TimeFormat))
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
