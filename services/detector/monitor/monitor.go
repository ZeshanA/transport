package monitor

import (
	"detector/fetch"
	"detector/request"
	"log"
	"time"
	"transport/lib/bus"
	"transport/lib/bustime"
	"transport/lib/database"
	"transport/lib/stringhelper"

	"github.com/VividCortex/ewma"
)

const (
	refreshInterval = 30 * time.Second
)

type TimedJourney struct {
	Duration                time.Duration
	DistanceFromArrivalTime time.Duration
}

func LiveBuses(avgTime int, predictedTime int, params request.JourneyParams, stopList []bustime.BusStop, complete chan bool) {
	ticker := time.NewTicker(refreshInterval)
	stopsAfterDest := stringhelper.SliceToSet(bustime.TrimStopList(stopList, params.ToStop, false))
	movingAverage := getInitialMovingAverage(avgTime, predictedTime)
	go monitorBuses(ticker, params, stopsAfterDest, complete, movingAverage)
	<-complete
}

func getInitialMovingAverage(avgTime int, predictedTime int) ewma.MovingAverage {
	movingAvg := ewma.NewMovingAverage()
	movingAvg.Add(float64(avgTime))
	movingAvg.Add(float64(predictedTime))
	return movingAvg
}

func monitorBuses(ticker *time.Ticker, params request.JourneyParams, stopsAfterDest map[string]bool, complete chan bool, movingAverageJourneyTime ewma.MovingAverage) {
	waitingToStart, startedJourneys := map[string]time.Time{}, map[string]time.Time{}
	for {
		select {
		case <-ticker.C:
			allJourneys, journeysApproachingStop := fetchJourneySets(params)
			findApproachingVehicles(journeysApproachingStop, waitingToStart)
			detectStartedJourneys(waitingToStart, startedJourneys, allJourneys, params)
			updateAverageJourneyTime(startedJourneys, allJourneys, stopsAfterDest, movingAverageJourneyTime)
			log.Println("Successfully processed live journey data, waiting for the next tick...")
		case <-complete:
			ticker.Stop()
			return
		}
	}
}

func updateAverageJourneyTime(startedJourneys map[string]time.Time, allJourneys map[string]bus.VehicleJourney,
	stopsAfterDest map[string]bool, movingAverageJourneyTime ewma.MovingAverage) {
	log.Println("Updating moving averages using any vehicles that have completed the route segment...")
	for vehicleRef, stamp := range startedJourneys {
		journey := allJourneys[vehicleRef]
		stopID := journey.StopPointRef.String
		if _, ok := stopsAfterDest[stopID]; ok {
			duration := journey.Timestamp.Sub(stamp).Seconds()
			log.Printf(
				"Vehicle with ID '%s' has completed its journey at %s, with a total duration of %d seconds",
				vehicleRef, stamp.Format(database.TimeFormat), duration,
			)
			movingAverageJourneyTime.Add(duration)
			delete(startedJourneys, vehicleRef)
			log.Printf("New average journey time: %f\n", movingAverageJourneyTime.Value())
		}
	}
}

func detectStartedJourneys(waitingToStart map[string]time.Time, journeyStarted map[string]time.Time, liveJourneys map[string]bus.VehicleJourney, params request.JourneyParams) {
	log.Println("Recording journey start times for any vehicles that have now moved onto our route segment...")
	for vehicleRef, stamp := range waitingToStart {
		if liveJourneys[vehicleRef].StopPointRef.String != params.FromStop {
			log.Printf("Vehicle with ID '%s' has started its journey at %s", vehicleRef, stamp.Format(database.TimeFormat))
			journeyStarted[vehicleRef] = stamp
			delete(waitingToStart, vehicleRef)
		}
	}
	log.Println("Successfully recorded journey start times")
}

func findApproachingVehicles(journeysApproachingStop map[string]bus.VehicleJourney, waitingToStart map[string]time.Time) {
	log.Println("Finding vehicles approaching the source stop...")
	for vehicleRef, journey := range journeysApproachingStop {
		waitingToStart[vehicleRef] = journey.Timestamp.Time
	}
	log.Printf("There are currently %d vehicles approaching the source stop\n", len(waitingToStart))
}

func fetchJourneySets(params request.JourneyParams) (liveJourneys map[string]bus.VehicleJourney, journeysApproachingStop map[string]bus.VehicleJourney) {
	log.Printf("Fetching live journeys for paramSet: %v...\n", params)
	liveJourneys = map[string]bus.VehicleJourney{}
	liveJourneys, err := fetch.LiveJourneys(params.RouteID, params.DirectionID)
	if err != nil {
		log.Fatalf("error fetching live journeys: %s", err)
	}
	journeysApproachingStop = extractJourneysApproachingStop(liveJourneys, params.FromStop)
	log.Println("Successfully fetched live journeys")
	return liveJourneys, journeysApproachingStop
}

func extractJourneysApproachingStop(journeys map[string]bus.VehicleJourney, stopID string) map[string]bus.VehicleJourney {
	matching := map[string]bus.VehicleJourney{}
	for vehicleID, journey := range journeys {
		if journey.StopPointRef.String == stopID {
			matching[vehicleID] = journey
		}
	}
	return matching
}
