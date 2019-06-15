package eval

import (
	"database/sql"
	"detector/calc"
	"detector/fetch"
	"detector/monitor"
	"detector/request"
	"detector/response"
	"fmt"
	"log"
	"math/rand"
	"sync"
	"time"
	"transport/lib/bus"
	"transport/lib/bustime"
	"transport/lib/database"
	"transport/lib/iohelper"
	"transport/lib/math"
	"transport/lib/nulltypes"
	"transport/lib/stringhelper"

	"github.com/VividCortex/ewma"
)

func Evaluate() {
	log.Println("Evaluation mode...")
	// Open a DB connection and schedule it to be closed after the program returns
	db := database.OpenDBConnection()
	defer db.Close()
	// Create a new BusTime client to handle metadata requests
	bt := bustime.NewClient(iohelper.GetEnv("MTA_API_KEY"))
	// Extract the number of journeys to evaluate
	numJourneys := 50
	var wg sync.WaitGroup
	wg.Add(numJourneys)
	log.Printf("Evaluating %d journeys...", numJourneys)
	for i := 0; i < numJourneys; i++ {
		params := generateRandomParams(bt)
		go performJourneyEvaluation(params, bt, db, wg)
	}
	wg.Wait()
}

func performJourneyEvaluation(params request.JourneyParams, bt *bustime.Client, db *sql.DB, wg sync.WaitGroup) {
	defer wg.Done()
	// Fetch the list of stops for the requested route and direction
	stops := bt.GetStops(params.RouteID)[params.RouteID][params.DirectionID]
	// Get average time to travel between stops
	avgTime, err := calc.AvgTimeBetweenStops(stops, params, db)
	if err != nil {
		log.Fatalf("error calculating average time between stops: %s", err)
	}
	log.Printf("Average time: %d\n", avgTime)
	log.Printf("Fetching predicted journey time...")
	predictedTime, err := fetch.PredictedJourneyTime(params, avgTime, stops)
	if err != nil {
		log.Fatalf("error calculating predicted time: %s", err)
	}
	log.Printf("Predicted time: %d\n", predictedTime)
	// Print current time and arrival time
	log.Printf("Time now is: %s", time.Now().In(database.TimeLoc).Format(database.TimeFormat))
	log.Printf("Arrival time is: %s", params.ArrivalTime.Format(database.TimeFormat))
	// Monitor live buses until we find a suitable vehicleID
	complete := make(chan response.Notification)
	monitor.LiveBuses(avgTime, predictedTime, params, stops, db, complete)
	notif := <-complete
	vehicleID := notif.VehicleID
	fmt.Printf("Suitable VehicleID found! Take the bus with ID %s\n", vehicleID)
	// Now time how long it actually takes the vehicleID to get to its stop
	ticker := time.NewTicker(30 * time.Second)
	timeTaken := make(chan float64)
	var startTime time.Time
	go func() {
		waitingToStart, startedJourneys := map[string]time.Time{}, map[string]time.Time{}
		stopsBeforeSource := stringhelper.SliceToSet(bustime.ExtractStops("before", params.FromStop, true, stops))
		stopsAfterDest := stringhelper.SliceToSet(bustime.ExtractStops("after", params.ToStop, false, stops))
		// Start with an empty average journey time
		journeyTime := ewma.NewMovingAverage()
		startedApproaching := false
		// Process what our vehicle is doing every time there's a new batch of live data
		for range ticker.C {
			// Get the set of all journeys, and a set of all journeys that are directly approaching the source stop
			allJourneys, _, journeysApproachingStop := monitor.FetchJourneySets(params, stopsBeforeSource)
			// Remove all approaching journeys that aren't related to the vehicle we're interested in
			if _, ok := journeysApproachingStop[vehicleID]; ok {
				log.Println("Found our vehicle in the 'Journeys Approaching Stop' list")
				journeysApproachingStop = map[string]bus.VehicleJourney{vehicleID: journeysApproachingStop[vehicleID]}
				startedApproaching = true
			} else {
				log.Println("Didn't find our vehicle in the 'Journeys Approaching Stop' list")
				journeysApproachingStop = map[string]bus.VehicleJourney{}
			}
			fmt.Println(allJourneys)
			// Store the vehicle's last-seen timestamp in `waitingToStart` if it still hasn't started its journey
			monitor.FindApproachingVehicles(journeysApproachingStop, waitingToStart)
			// At the timestamp when the vehicle starts its journey, add the *previous* timestamp
			// (its 'start' timestamp) to startedJourneys
			monitor.DetectStartedJourneys(waitingToStart, startedJourneys, allJourneys, params)
			// If the vehicle starts to approach a stop that is later on the route than our destination stop,
			// mark its journey as completed and record the time the journey took
			monitor.UpdateAverageJourneyTime(startedJourneys, allJourneys, stopsAfterDest, journeyTime)
			// Is the vehicle waiting to start the journey, or has it started?
			_, waiting := waitingToStart[vehicleID]
			startStamp, started := startedJourneys[vehicleID]
			if started {
				startTime = startStamp
			}
			if startedApproaching && !waiting && !started {
				timeTaken <- journeyTime.Value()
			}
		}
	}()
	actualJourneyTime := <-timeTaken
	expectedAt := params.ArrivalTime.Time
	arrivedAt := startTime.Add(time.Duration(actualJourneyTime) * time.Second)
	offBy := int(arrivedAt.Sub(expectedAt) / time.Second)
	fmt.Printf(
		"StartTime: %s\nExpected At: %s\nArrived At: %s\nOff by %d",
		startTime,
		expectedAt.Format(database.TimeFormat),
		arrivedAt.Format(database.TimeFormat),
		offBy,
	)
	entries := []NotificationEval{{
		params.RouteID, params.DirectionID,
		params.FromStop, params.ToStop,
		nulltypes.TimestampFrom(params.ArrivalTime),
		nulltypes.TimestampFrom(database.Timestamp{Time: arrivedAt}),
		offBy,
	}}
	database.Store(database.NotificationEvalTable, func(evalEntry interface{}) []interface{} {
		entry := evalEntry.(NotificationEval)
		return []interface{}{
			entry.RouteID, entry.DirectionID,
			entry.FromStop, entry.ToStop,
			entry.DesiredArrivalTime, entry.ActualArrivalTime,
			entry.OffBy,
		}
	}, NotificationEvalToInterface(entries))
}

type NotificationEval struct {
	RouteID            string
	DirectionID        int
	FromStop           string
	ToStop             string
	DesiredArrivalTime nulltypes.Timestamp
	ActualArrivalTime  nulltypes.Timestamp
	OffBy              int
}

func NotificationEvalToInterface(entries []NotificationEval) []interface{} {
	r := make([]interface{}, len(entries))
	for i, journey := range entries {
		r[i] = journey
	}
	return r
}

func generateRandomParams(bt *bustime.Client) request.JourneyParams {
	log.Println("Generating random parameter set...")
	rj, err := fetch.RawJourneys()
	if err != nil {
		log.Fatalf("error fetching raw journeys: %v", err)
	}
	rjArr := rj.Array()
	params := request.JourneyParams{}
	var stops []bustime.BusStop
	for {
		// Select a random journey from the list of currently active ones
		randomMvmt := rjArr[rand.Intn(len(rjArr))]
		routeID := randomMvmt.Get("LineRef").String()
		if !validRouteID(routeID) {
			continue
		}
		directionID := int(randomMvmt.Get("DirectionRef").Int())
		nextStop := randomMvmt.Get("StopPointRef").String()
		stops = bt.GetStops(routeID)[routeID][directionID]
		isLastStop := stops[len(stops)-1].ID == nextStop
		// If we're going to the last stop, no further journeys are possible - try again
		if isLastStop {
			continue
		}
		params.RouteID, params.DirectionID, params.FromStop = routeID, directionID, nextStop
		break
	}
	// Pick random destination stop
	stopsAfterSource := bustime.ExtractStops("after", params.FromStop, false, stops)
	destStop := stops[rand.Intn(math.MaxInt(len(stopsAfterSource)-1, 1))]
	params.ToStop = destStop.ID
	// Pick random arrival time
	delay := time.Duration(math.RandInRange(15, 30)) * time.Minute
	params.ArrivalTime = database.Timestamp{Time: time.Now().In(database.TimeLoc).Add(delay)}
	log.Printf("Selected random parameter set: %s", params.String())
	return params
}

func validRouteID(routeID string) bool {
	validRoutes := []string{
		"MTA NYCT_M102", "MTA NYCT_S86", "MTA NYCT_SIM8X", "MTA NYCT_SIM4X", "MTABC_QM36", "MTABC_QM44", "MTABC_QM31",
		"MTABC_QM40", "MTABC_QM42", "MTABC_QM34", "MTA NYCT_SIM32", "MTA NYCT_S81", "MTA NYCT_SIM35", "MTABC_QM32",
		"MTA NYCT_SIM9", "MTABC_QM35"}
	for _, id := range validRoutes {
		if id == routeID {
			return true
		}
	}
	return false
}
