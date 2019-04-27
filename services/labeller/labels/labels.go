package labels

import (
	"fmt"
	"labeller/stopdistance"
	"log"
	"time"
	"transport/lib/bus"
)

func Create(partitionedJourneys map[bus.DirectedRoute][]bus.VehicleJourney, stopDistances map[stopdistance.Key]float64) (labelledMvmts []bus.LabelledJourney) {
	for route, mvmts := range partitionedJourneys {
		x, y, z := 0, 0, 0
		if len(mvmts) < 2 {
			log.Println("Fewer than 2 movements for route.")
			continue
		}
		for i, fromMvmt := range mvmts {
			finalPreStopMvmt, reachedStopMvmt, wentPastStop := ExtractKeyMvmts(mvmts, i)
			// No movement event where the vehicle reaches its stop or moves to the next stopID: we can't ascertain
			// arrival time, so skip the current movement event.
			// TODO: could potentially extrapolate it using the same technique as wentPastStop
			if reachedStopMvmt == nil {
				x += 1
				continue
			} else if !wentPastStop {
				// Perfect stop, so we can just subtract the timestamps to calculate how long the journey took
				y += 1
				timeTaken := SubtractMvmtTimestamps(reachedStopMvmt, &fromMvmt)
				labelledMvmts = append(labelledMvmts, bus.LabelledJourneyFrom(fromMvmt, int(timeTaken)))
			} else {
				// Went past stop: need to extrapolate time
				z += 1
				timeToFinalPreStopMvmt := SubtractMvmtTimestamps(finalPreStopMvmt, &fromMvmt)
				timeFromFinalPreStopMvmtToStop := GetTimeToStopFromFinalMovement(route, fromMvmt, finalPreStopMvmt, reachedStopMvmt, stopDistances)
				totalTimeFromMvmtToStop := float64(timeToFinalPreStopMvmt) + timeFromFinalPreStopMvmtToStop
				labelledMvmts = append(labelledMvmts, bus.LabelledJourneyFrom(fromMvmt, int(totalTimeFromMvmtToStop)))
			}
			fmt.Printf("No stop: %d, Exact stop: %d, Went Past Stop: %d, Total: %d\n", x, y, z, x+y+z)
		}
	}
	fmt.Println(count)
	fmt.Println(total)
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

func GetTimeToStopFromFinalMovement(route bus.DirectedRoute, fromMvmt bus.VehicleJourney, preStopMvmt *bus.VehicleJourney, postStopMvmt *bus.VehicleJourney, stopDistances map[stopdistance.Key]float64) float64 {
	distanceBetweenStops := GetDistanceBetweenStops(route, fromMvmt, postStopMvmt, stopDistances)
	distanceToNextStop := float64(postStopMvmt.DistanceFromStop.Int64)
	distanceTravelledPastPrevStop := distanceBetweenStops - distanceToNextStop
	distanceFromFinalPreStopMvmtToStop := float64(preStopMvmt.DistanceFromStop.Int64)
	distanceTravelledBetweenMvmts := distanceFromFinalPreStopMvmtToStop + distanceTravelledPastPrevStop
	timeBetweenMvmtEvents := SubtractMvmtTimestamps(postStopMvmt, preStopMvmt)
	speedBetweenMvmtEvents := distanceTravelledBetweenMvmts / timeBetweenMvmtEvents
	timeFromFinalPreStopMvmtToStop := distanceFromFinalPreStopMvmtToStop / speedBetweenMvmtEvents
	return timeFromFinalPreStopMvmtToStop
}

var count = 0
var total = 0

func GetDistanceBetweenStops(route bus.DirectedRoute, from bus.VehicleJourney, to *bus.VehicleJourney, stopDistances map[stopdistance.Key]float64) float64 {
	prevStopID, nextStopID := from.StopPointRef.String, to.StopPointRef.String
	key := stopdistance.Key{RouteID: route.RouteID, DirectionID: route.DirectionID, FromID: prevStopID, ToID: nextStopID}
	if _, ok := stopDistances[key]; true {
		fmt.Println(key)
		if !ok {
			count += 1
		}
		total += 1
	}
	return stopDistances[key]
}

func SubtractMvmtTimestamps(mvmtA, mvmtB *bus.VehicleJourney) float64 {
	timeDiff := mvmtA.Timestamp.Timestamp.Sub(mvmtB.Timestamp.Timestamp.Time)
	return float64(timeDiff) / float64(time.Second)
}
