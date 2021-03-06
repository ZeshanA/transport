package monitor

import (
	"database/sql"
	"detector/calc"
	"detector/fetch"
	"detector/request"
	"detector/response"
	"fmt"
	"log"
	"time"
	"transport/lib/bus"
	"transport/lib/bustime"
	"transport/lib/database"
	"transport/lib/stringhelper"

	"github.com/VividCortex/ewma"
)

const (
	RefreshInterval = 30 * time.Second
)

type TimedJourney struct {
	Duration                time.Duration
	DistanceFromArrivalTime time.Duration
}

func LiveBuses(avgTime int, predictedTime int, params request.JourneyParams, stopList []bustime.BusStop, db *sql.DB, complete chan response.Notification) {
	ticker := time.NewTicker(RefreshInterval)
	movingAverage := GetInitialMovingAverage(avgTime, predictedTime)
	go monitorBuses(ticker, params, stopList, db, complete, movingAverage)
	<-complete
}

func GetInitialMovingAverage(avgTime int, predictedTime int) ewma.MovingAverage {
	movingAvg := ewma.NewMovingAverage()
	movingAvg.Add(float64(avgTime))
	movingAvg.Add(float64(predictedTime))
	return movingAvg
}

func monitorBuses(ticker *time.Ticker, params request.JourneyParams, stopList []bustime.BusStop, db *sql.DB, complete chan response.Notification, movingAvg ewma.MovingAverage) {
	waitingToStart, startedJourneys := map[string]time.Time{}, map[string]time.Time{}
	stopsBeforeSource := stringhelper.SliceToSet(bustime.ExtractStops("before", params.FromStop, true, stopList))
	stopsAfterDest := stringhelper.SliceToSet(bustime.ExtractStops("after", params.ToStop, false, stopList))
	for {
		select {
		case <-ticker.C:
			allJourneys, journeysBeforeStop, journeysApproachingStop := FetchJourneySets(params, stopsBeforeSource)
			FindApproachingVehicles(journeysApproachingStop, waitingToStart)
			DetectStartedJourneys(waitingToStart, startedJourneys, allJourneys, params)
			UpdateAverageJourneyTime(startedJourneys, allJourneys, stopsAfterDest, movingAvg)
			idealArrivalTime := params.ArrivalTime.Add(-time.Duration(movingAvg.Value()) * time.Second)
			for vehicleID, journey := range journeysBeforeStop {
				timeToNextStop, err := fetch.SingleMovementPrediction(journey)
				if err != nil {
					log.Printf("error fetching predicted journey time: %s", err)
				}
				nxtStopToSourceStopParams := params
				nxtStopToSourceStopParams.FromStop = journey.StopPointRef.String
				nxtStopToSourceStopParams.ToStop = params.FromStop
				nxtStopToSourceStopParams.ArrivalTime = database.Timestamp{Time: idealArrivalTime}
				avgTime, err := calc.AvgTimeBetweenStops(stopList, nxtStopToSourceStopParams, db)
				if err != nil {
					log.Fatalf("error calculating average time between stops: %s", err)
				}
				nextStopToSourceStop, err := fetch.PredictedJourneyTime(nxtStopToSourceStopParams, avgTime, stopList)
				totalTimeUntilSourceStop := timeToNextStop + nextStopToSourceStop
				currentVehicleArrivalTime := time.Now().In(database.TimeLoc).Add(time.Duration(totalTimeUntilSourceStop) * time.Second)
				diff := idealArrivalTime.Sub(currentVehicleArrivalTime)
				log.Printf("We would like a bus arriving at exactly %s", idealArrivalTime)
				log.Printf("Vehicle with ID %s is estimated to arrive at the source stop in %d seconds, at %s", vehicleID, totalTimeUntilSourceStop, currentVehicleArrivalTime.Format(database.TimeFormat))
				log.Printf("The gap between these two times is %f seconds", diff.Seconds())
				if diff < 5*time.Minute {
					complete <- response.Notification{
						VehicleID:            vehicleID,
						OptimalDepartureTime: database.Timestamp{Time: currentVehicleArrivalTime},
						PredictedArrivalTime: database.Timestamp{Time: idealArrivalTime.Add(time.Duration(movingAvg.Value()) * time.Second)},
					}
				}
			}
			log.Println("Successfully processed live journey data, waiting for the next tick...")
		case <-complete:
			ticker.Stop()
			return
		}
	}
}

func UpdateAverageJourneyTime(startedJourneys map[string]time.Time, allJourneys map[string]bus.VehicleJourney,
	stopsAfterDest map[string]bool, movingAverageJourneyTime ewma.MovingAverage) {
	log.Println("Updating moving averages using any vehicles that have completed the route segment...")
	for vehicleRef, startTimestamp := range startedJourneys {
		journey := allJourneys[vehicleRef]
		stopID := journey.StopPointRef.String
		if _, ok := stopsAfterDest[stopID]; ok {
			duration := journey.Timestamp.Sub(startTimestamp).Seconds()
			log.Printf(
				"Vehicle with ID '%s' has completed its journey at %s, with a total duration of %f seconds",
				vehicleRef, journey.Timestamp.Format(database.TimeFormat), duration,
			)
			movingAverageJourneyTime.Add(duration)
			delete(startedJourneys, vehicleRef)
			log.Printf("New average journey time: %f\n", movingAverageJourneyTime.Value())
		}
	}
}

func DetectStartedJourneys(waitingToStart map[string]time.Time, startedJourneys map[string]time.Time, liveJourneys map[string]bus.VehicleJourney, params request.JourneyParams) {
	log.Println("Recording journey start times for any vehicles that have now moved onto our route segment...")
	for vehicleRef, stamp := range waitingToStart {
		if liveJourneys[vehicleRef].StopPointRef.String != params.FromStop {
			log.Printf("Vehicle with ID '%s' has started its journey at %s", vehicleRef, stamp.Format(database.TimeFormat))
			startedJourneys[vehicleRef] = stamp
			delete(waitingToStart, vehicleRef)
		}
	}
	log.Println("Successfully recorded journey start times")
}

func FindApproachingVehicles(journeysApproachingStop map[string]bus.VehicleJourney, waitingToStart map[string]time.Time) {
	log.Println("Finding vehicles approaching the source stop...")
	for vehicleRef, journey := range journeysApproachingStop {
		fmt.Println("SETTING STAMP TO:")
		fmt.Println(journey.Timestamp.Time.Format(database.TimeFormat))
		waitingToStart[vehicleRef] = journey.Timestamp.Time
	}
	log.Printf("There are currently %d vehicles approaching the source stop\n", len(waitingToStart))
}

func FetchJourneySets(params request.JourneyParams, stopsBeforeStopID map[string]bool) (allJourneys map[string]bus.VehicleJourney, journeysBeforeStop map[string]bus.VehicleJourney, journeysApproachingStop map[string]bus.VehicleJourney) {
	log.Printf("Fetching live journeys for paramSet: %v...\n", params)
	allJourneys = map[string]bus.VehicleJourney{}
	allJourneys, err := fetch.LiveJourneys(params.RouteID, params.DirectionID)
	if err != nil {
		log.Fatalf("error fetching live journeys: %s", err)
	}
	journeysBeforeStop = extractJourneysBeforeStop(allJourneys, stopsBeforeStopID)
	journeysApproachingStop = extractJourneysApproachingStop(allJourneys, params.FromStop)
	log.Println("Successfully fetched live journeys")
	return allJourneys, journeysBeforeStop, journeysApproachingStop
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

func extractJourneysBeforeStop(journeys map[string]bus.VehicleJourney, stopsBeforeStopID map[string]bool) map[string]bus.VehicleJourney {
	matching := map[string]bus.VehicleJourney{}
	for vehicleID, journey := range journeys {
		if _, exists := stopsBeforeStopID[journey.StopPointRef.String]; exists {
			matching[vehicleID] = journey
		}
	}
	return matching
}
