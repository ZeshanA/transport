package calc

import (
	"database/sql"
	"detector/fetch"
	"detector/request"
	"transport/lib/bus"
	"transport/lib/bustime"
	"transport/lib/math"
)

type Journey struct {
	PreStop  bus.LabelledJourney
	PostStop bus.LabelledJourney
}

// Get the average time taken for vehicles to travel between the two stops
// around the requested arrival time
func AvgTimeBetweenStops(stopList []bustime.BusStop, jp request.JourneyParams, db *sql.DB) (int, error) {
	// Fetch movements that match the requested parameters
	mvmts, err := fetch.MovementsInWindow(db, stopList, jp)
	if err != nil {
		return 0, nil
	}
	// Split movements up by vehicleID
	splitMvmts := SplitMovementsByVehicleID(mvmts)
	// Get a list containing how long each 'fromStop' -> 'toStop' journey took in seconds
	journeyTimes := mvmtsToJourneyTimes(splitMvmts, stopList, jp)
	// TODO: Switch to median
	avgJourneyTime := math.SliceMeanRounded(journeyTimes)
	return avgJourneyTime, nil
}

// Takes a list of movements (split by vehicleID) and returns a list containing the durations of individual
// journeys between the two stops in the request.JourneyParams struct
func mvmtsToJourneyTimes(splitMvmts map[string][]bus.LabelledJourney, stopList []bustime.BusStop, jp request.JourneyParams) []int {
	var journeyTimes []int
	for _, mvmts := range splitMvmts {
		journeys := getPreAndPostStopMvmts(mvmts, stopList, jp)
		for _, journey := range journeys {
			preStop, postStop := journey.PreStop.Timestamp, journey.PostStop.Timestamp
			journeyTimes = append(journeyTimes, int(postStop.Sub(preStop.Time).Seconds()))
		}
	}
	return journeyTimes
}

func getPreAndPostStopMvmts(mvmts []bus.LabelledJourney, stopList []bustime.BusStop, jp request.JourneyParams) (journeys []Journey) {
	var curJourney Journey
	// If we see any of these stops, we know we've gone past the destination stop (i.e. the journey has ended)
	stopsAfterDestinationStop := stopsAfter(jp.ToStop, stopList)
	// Flags to indicate which state the search is in
	lookingForJourneyEnd, foundFinalPreStop := false, false
	// Loop over all movements
	for _, mvmt := range mvmts {
		// Waiting for a movement that starts the journey (i.e. is approaching the 'fromStop')
		if !lookingForJourneyEnd {
			// Skip until we find a movement approaching the fromStop
			if mvmt.StopPointRef.String != jp.FromStop {
				continue
			} else {
				// Found a movement approaching the 'fromStop'
				lookingForJourneyEnd = true
				curJourney.PreStop = mvmt
			}
		} else {
			// We currently already have a previous movement that *could* start a journey,
			// but we don't know if it's the final pre-stop movement

			// If the current movement's stopID isn't the 'fromStop', then we've gone past the 'fromStop
			// and we now know that the previous movement we've stored is the final pre-stop movement.
			if !foundFinalPreStop && mvmt.StopPointRef.String != jp.FromStop {
				foundFinalPreStop = true
			}

			// We now know that we've found the final pre-stop movement and need to find
			// the first movement where the bus has gone past the 'toStop' (i.e. the end of the journey)
			if foundFinalPreStop {
				// Found a post-stop movement
				if _, pastDestStop := stopsAfterDestinationStop[mvmt.StopPointRef.String]; pastDestStop {
					curJourney.PostStop = mvmt
					journeys = append(journeys, curJourney)
					// Reset flags and current journey
					curJourney, lookingForJourneyEnd, foundFinalPreStop = Journey{}, false, false
				}
			}
		}
	}
	return journeys
}

// Returns a map (mimicking a Set) containing the IDs of all
// stops that come *after* the given stopID on the route.
func stopsAfter(stopID string, stopList []bustime.BusStop) map[string]bool {
	result := make(map[string]bool)
	stopFound := false
	for _, stop := range stopList {
		if stop.ID == stopID {
			stopFound = true
			continue
		}
		if stopFound {
			result[stop.ID] = true
		}
	}
	return result
}

// Takes a list of movements and returns those same movements as a map, keyed by vehicleID
func SplitMovementsByVehicleID(mvmts []bus.LabelledJourney) map[string][]bus.LabelledJourney {
	splitMvmts := make(map[string][]bus.LabelledJourney)
	for _, mvmt := range mvmts {
		splitMvmts[mvmt.VehicleRef.String] = append(splitMvmts[mvmt.VehicleRef.String], mvmt)
	}
	return splitMvmts
}
