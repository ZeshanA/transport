package labels

import (
	"labeller/stopdistance"
	"time"
	"transport/lib/bus"
)

// Create takes a map of journeys partitioned by their DirectedRoute
// returns: a slice of labelledMovements that can be inserted into the DB.
func Create(partitionedJourneys map[bus.DirectedRoute][]bus.VehicleJourney, stopDistances map[stopdistance.Key]float64, averageStopDistances map[string]int) (labelledMvmts []bus.LabelledJourney) {
	for route, mvmts := range partitionedJourneys {
		if len(mvmts) < 2 {
			continue
		}
		labelledMvmts = append(labelledMvmts, labelMvmtsForRoute(route, mvmts, stopDistances, averageStopDistances)...)
	}
	return labelledMvmts
}

// Returns a slice of labelledJourneys for a single route
func labelMvmtsForRoute(route bus.DirectedRoute, mvmts []bus.VehicleJourney, stopDistances map[stopdistance.Key]float64, averageStopDistances map[string]int) []bus.LabelledJourney {
	var labelledMvmts []bus.LabelledJourney
	for i, fromMvmt := range mvmts {
		finalPreStopMvmt, reachedStopMvmt, wentPastStop := ExtractKeyMvmts(mvmts, i)
		// No movement event where the vehicle reaches its stop or moves to the next stopID: we can't ascertain
		// arrival time, so skip the current movement event.
		// TODO: could potentially extrapolate it using the same technique as wentPastStop
		if reachedStopMvmt == nil {
			continue
		} else if !wentPastStop {
			// Perfect stop, so we can just subtract the timestamps to calculate how long the journey took
			timeTaken := SubtractMvmtTimestamps(reachedStopMvmt, &fromMvmt)
			labelledMvmts = append(labelledMvmts, bus.LabelledJourneyFrom(fromMvmt, int(timeTaken)))
		} else {
			// Went past stop: need to extrapolate time
			timeToFinalPreStopMvmt := SubtractMvmtTimestamps(finalPreStopMvmt, &fromMvmt)
			timeFromFinalPreStopMvmtToStop := GetTimeToStopFromFinalMovement(route, fromMvmt, finalPreStopMvmt, reachedStopMvmt, stopDistances, averageStopDistances)
			totalTimeFromMvmtToStop := float64(timeToFinalPreStopMvmt) + timeFromFinalPreStopMvmtToStop
			labelledMvmts = append(labelledMvmts, bus.LabelledJourneyFrom(fromMvmt, int(totalTimeFromMvmtToStop)))
		}
	}
	return labelledMvmts
}

// finalPreStopMvmt: The last movement before the vehicle changed stopIDs. (only relevant for
// extrapolating times when the bus goes past a stop without sending a movement event)
// reachedStopMvmt: The first movement at which the vehicle either reached the next stop
// or went past it.
// wentPastStop: Did the vehicle go past the stop without a movementEvent with distanceFromStop = 0?
// If the vehicle goes past a stop and then sends a movement event, the stopID will change,
// but there will not be a movement event telling us it "reached" the last stop (startIndex.e. there won't
// be an event with distanceFromStop = 0 for the previous stop).
func ExtractKeyMvmts(mvmts []bus.VehicleJourney, startIndex int) (finalPreStopMvmt *bus.VehicleJourney, reachedStopMvmt *bus.VehicleJourney, wentPastStop bool) {
	startingAtAStop := mvmts[startIndex].DistanceFromStop.Int64 == 0
	wentPastStop = false
	for j := startIndex + 1; j < len(mvmts); j++ {
		if !startingAtAStop && mvmts[startIndex].StopPointRef == mvmts[j].StopPointRef && mvmts[j].DistanceFromStop.Int64 == 0 {
			reachedStopMvmt, finalPreStopMvmt = &mvmts[j], &mvmts[j]
			break
		}
		if mvmts[j].StopPointRef.String != mvmts[startIndex].StopPointRef.String {
			reachedStopMvmt, finalPreStopMvmt = &mvmts[j], &mvmts[j-1]
			if !startingAtAStop {
				wentPastStop = true
			}
			break
		}
	}
	return finalPreStopMvmt, reachedStopMvmt, wentPastStop
}

// Returns an estimate of how long it took the bus to reach the stop from the position it was at when
// we received the final pre-stop movement event.
func GetTimeToStopFromFinalMovement(route bus.DirectedRoute, fromMvmt bus.VehicleJourney, preStopMvmt *bus.VehicleJourney, postStopMvmt *bus.VehicleJourney, stopDistances map[stopdistance.Key]float64, averageStopDistances map[string]int) float64 {
	distanceBetweenStops := GetDistanceBetweenStops(route, fromMvmt, postStopMvmt, stopDistances, averageStopDistances)
	distanceToNextStop := float64(postStopMvmt.DistanceFromStop.Int64)
	distanceTravelledPastPrevStop := distanceBetweenStops - distanceToNextStop
	if distanceTravelledPastPrevStop < 0 {
		distanceTravelledPastPrevStop = bus.AverageDistanceBetweenStops * 0.1
	}
	distanceFromFinalPreStopMvmtToStop := float64(preStopMvmt.DistanceFromStop.Int64)
	distanceTravelledBetweenMvmts := distanceFromFinalPreStopMvmtToStop + distanceTravelledPastPrevStop
	timeBetweenMvmtEvents := SubtractMvmtTimestamps(postStopMvmt, preStopMvmt)
	speedBetweenMvmtEvents := distanceTravelledBetweenMvmts / timeBetweenMvmtEvents
	timeFromFinalPreStopMvmtToStop := distanceFromFinalPreStopMvmtToStop / speedBetweenMvmtEvents
	return timeFromFinalPreStopMvmtToStop
}

func GetDistanceBetweenStops(route bus.DirectedRoute, from bus.VehicleJourney, to *bus.VehicleJourney, stopDistances map[stopdistance.Key]float64, averageStopDistances map[string]int) float64 {
	prevStopID, nextStopID := from.StopPointRef.String, to.StopPointRef.String
	key := stopdistance.Key{RouteID: route.RouteID, DirectionID: route.DirectionID, FromID: prevStopID, ToID: nextStopID}
	if preciseDist, found := stopDistances[key]; found {
		// Found precise distance between the two stops
		return preciseDist
	} else if avgDist, found := averageStopDistances[bus.RemoveAgencyID(route.RouteID)]; found {
		// Couldn't find the precise distance between the stops, use average distance for stops on that route.
		return float64(avgDist)
	} else {
		// Couldn't find the average distance for this route, use average distance for stops across NYC.
		return bus.AverageDistanceBetweenStops
	}
}

func SubtractMvmtTimestamps(mvmtA, mvmtB *bus.VehicleJourney) float64 {
	timeDiff := mvmtA.Timestamp.Timestamp.Sub(mvmtB.Timestamp.Timestamp.Time)
	return float64(timeDiff) / float64(time.Second)
}
