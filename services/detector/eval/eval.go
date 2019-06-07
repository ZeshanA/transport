package eval

import (
	"detector/calc"
	"detector/fetch"
	"detector/monitor"
	"detector/request"
	"fmt"
	"log"
	"math/rand"
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
	log.Println("Evaluation mode")
	// Open a DB connection and schedule it to be closed after the program returns
	db := database.OpenDBConnection()
	defer db.Close()
	bt := bustime.NewClient(iohelper.GetEnv("MTA_API_KEY"))
	// Generate random set of parameters
	params := generateRandomParams(bt)

	// Fetch the list of stops for the requested route and direction
	stops := bt.GetStops(params.RouteID)[params.RouteID][params.DirectionID]
	fmt.Println(stops)
	// Get average time to travel between stops
	avgTime, err := calc.AvgTimeBetweenStops(stops, params, db)
	if err != nil {
		log.Fatalf("error calculating average time between stops: %s", err)
	}
	log.Printf("Average time: %d\n", avgTime)
	predictedTime, err := fetch.PredictedJourneyTime(params, avgTime, stops)
	if err != nil {
		log.Fatalf("error calculating predicted time: %s", err)
	}
	log.Printf("Predicted time: %d\n", predictedTime)
	log.Printf("Time now is: %s", time.Now().In(database.TimeLoc).Format(database.TimeFormat))
	log.Printf("Arrival time is: %s", params.ArrivalTime.In(database.TimeLoc).Format(database.TimeFormat))
	complete := make(chan string)
	monitor.LiveBuses(avgTime, predictedTime, params, stops, db, complete)
	vehicleID := <-complete
	fmt.Println(vehicleID)
	// Now time how long it takes the vehicleID to get to its stop
	ticker := time.NewTicker(1 * time.Second)
	timeTaken := make(chan float64)

	waitingToStart, startedJourneys := map[string]time.Time{}, map[string]time.Time{}
	stopsBeforeSource := stringhelper.SliceToSet(bustime.ExtractStops("before", params.FromStop, true, stops))
	stopsAfterDest := stringhelper.SliceToSet(bustime.ExtractStops("after", params.ToStop, false, stops))
	var startTime time.Time
	go func() {
		journeyTime := ewma.NewMovingAverage()
		for range ticker.C {
			allJourneys, _, journeysApproachingStop := monitor.FetchJourneySets(params, stopsBeforeSource)
			if _, ok := journeysApproachingStop[vehicleID]; ok {
				journeysApproachingStop = map[string]bus.VehicleJourney{vehicleID: journeysApproachingStop[vehicleID]}
			} else {
				journeysApproachingStop = map[string]bus.VehicleJourney{}
			}
			fmt.Println(allJourneys)
			monitor.FindApproachingVehicles(journeysApproachingStop, waitingToStart)
			monitor.DetectStartedJourneys(waitingToStart, startedJourneys, allJourneys, params)
			monitor.UpdateAverageJourneyTime(startedJourneys, allJourneys, stopsAfterDest, journeyTime)
			_, waiting := waitingToStart[vehicleID]
			startStamp, started := startedJourneys[vehicleID]
			if started {
				startTime = startStamp
			}
			if !waiting && !started {
				timeTaken <- journeyTime.Value()
			}
		}
	}()
	actualJourneyTime := <-timeTaken
	expectedAt := params.ArrivalTime
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
		nulltypes.TimestampFrom(database.Timestamp{Time: params.ArrivalTime}),
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
	// Get list of all routeIDs
	routeIDs := bt.GetRoutes("MTA NYCT", "MTABC")
	// Initialise global pseudo-RNG
	rand.Seed(time.Now().Unix())
	// Select a random routeID and directionID
	randomRouteID := routeIDs[rand.Intn(len(routeIDs))]
	possibleDirections := 2
	randomDirectionID := rand.Intn(possibleDirections)
	// Fetch stops for the given routeID and directionID
	stops := bt.GetStops(randomRouteID)[randomRouteID][randomDirectionID]
	// Select a random source stop (excluding the final stop on the route, as there are
	// no possible journeys in that case)
	sourceStopIndex := rand.Intn(len(stops) - 1)
	sourceStop := stops[sourceStopIndex].ID
	destStop := stops[math.RandInRange(sourceStopIndex+1, len(stops))].ID
	// Pick random arrival time
	delay := time.Duration(math.RandInRange(1, 5)) * time.Hour
	arrivalTime := time.Now().In(database.TimeLoc).Add(delay)
	return request.JourneyParams{
		RouteID: randomRouteID, DirectionID: randomDirectionID,
		FromStop: sourceStop, ToStop: destStop,
		ArrivalTime: arrivalTime,
	}
}
