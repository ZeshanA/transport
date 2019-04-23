package main

import (
	"fmt"
	"time"
	"transport/lib/bus"
	"transport/lib/database"
	"transport/lib/nulltypes"

	"gopkg.in/guregu/null.v3"
)

var dbConn = database.OpenDBConnection()
var dateRange = []time.Time{
	time.Date(2019, 04, 21, 0, 0, 0, 0, time.UTC),
	time.Date(2019, 04, 22, 0, 0, 0, 0, time.UTC),
}

//func main() {
//	stopDistances := stopdistance.Get(dbConn)
//	fmt.Println(stopDistances)
//	dataForDates := fetch.DateRange(dbConn, dateRange[0], dateRange[1])
//	partitionedDataForDates := make([]map[bus.DirectedRoute][]bus.VehicleJourney, len(dataForDates))
//	for i, journeysOnDate := range dataForDates {
//		partitionedDataForDates[i] = bus.PartitionJourneys(journeysOnDate)
//	}
//	fmt.Println(partitionedDataForDates)
//}

func main() {
	reachesStop := map[bus.DirectedRoute][]bus.VehicleJourney{
		{"M55", 0, "ABC"}: {
			{
				DistanceFromStop: null.IntFrom(200), StopPointRef: null.StringFrom("1"),
				Timestamp: nulltypes.TimestampFrom(database.Timestamp{Time: time.Date(2019, 04, 23, 16, 35, 00, 0, time.UTC)}),
			},
			{
				DistanceFromStop: null.IntFrom(100), StopPointRef: null.StringFrom("1"),
				Timestamp: nulltypes.TimestampFrom(database.Timestamp{Time: time.Date(2019, 04, 23, 16, 34, 00, 0, time.UTC)}),
			},
			{
				DistanceFromStop: null.IntFrom(50), StopPointRef: null.StringFrom("1"),
				Timestamp: nulltypes.TimestampFrom(database.Timestamp{Time: time.Date(2019, 04, 23, 16, 32, 00, 0, time.UTC)}),
			},
			{
				DistanceFromStop: null.IntFrom(0), StopPointRef: null.StringFrom("1"),
				Timestamp: nulltypes.TimestampFrom(database.Timestamp{Time: time.Date(2019, 04, 23, 16, 30, 00, 0, time.UTC)}),
			},
			{
				DistanceFromStop: null.IntFrom(0), StopPointRef: null.StringFrom("1"),
				Timestamp: nulltypes.TimestampFrom(database.Timestamp{Time: time.Date(2019, 04, 23, 16, 29, 00, 0, time.UTC)}),
			},
		},
	}
	CreateLabels(reachesStop)
	fmt.Println("DOESN'T REACH STOP:::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::")

	stopNotReached := map[bus.DirectedRoute][]bus.VehicleJourney{
		{"M55", 0, "ABC"}: {
			{
				DistanceFromStop: null.IntFrom(200), StopPointRef: null.StringFrom("1"),
				Timestamp: nulltypes.TimestampFrom(database.Timestamp{Time: time.Date(2019, 04, 23, 16, 35, 00, 0, time.UTC)}),
			},
			{
				DistanceFromStop: null.IntFrom(100), StopPointRef: null.StringFrom("1"),
				Timestamp: nulltypes.TimestampFrom(database.Timestamp{Time: time.Date(2019, 04, 23, 16, 34, 00, 0, time.UTC)}),
			},
			{
				DistanceFromStop: null.IntFrom(50), StopPointRef: null.StringFrom("1"),
				Timestamp: nulltypes.TimestampFrom(database.Timestamp{Time: time.Date(2019, 04, 23, 16, 32, 00, 0, time.UTC)}),
			},
		},
	}
	CreateLabels(stopNotReached)
}

func CreateLabels(partitionedJourneys map[bus.DirectedRoute][]bus.VehicleJourney) {
	for _, mvmts := range partitionedJourneys {
		if len(mvmts) < 2 {
			break
		}
		for i, fromMvmt := range mvmts {
			if fromMvmt.DistanceFromStop.Int64 == 0 {
				continue
			}
			// The first movement at which the vehicle either reached the next stop
			// or went past it.
			var reachedStopMvmt *bus.VehicleJourney
			// Did the vehicle go past the stop without a movementEvent with distanceFromStop = 0?
			// If the vehicle goes past a stop and then sends a movement event, the stopID will change,
			// but there will not be a movement event telling us it "reached" the last stop (i.e. there won't
			// be an event with distanceFromStop = 0 for the previous stop).
			wentPastStop := false
			for j := i + 1; j < len(mvmts); j++ {
				if mvmts[j].DistanceFromStop.Int64 == 0 {
					reachedStopMvmt = &mvmts[j]
					break
				}
				if mvmts[j].StopPointRef.String != fromMvmt.StopPointRef.String {
					reachedStopMvmt = &mvmts[j]
					wentPastStop = true
					break
				}
			}
			// If there is no movement event where the vehicle reaches its stop
			// or moves to the next stopID, then we can't ascertain arrival time,
			// so we skip the current movement event.
			// TODO: could potentially extrapolate it using the same technique as wentPastStop
			if reachedStopMvmt == nil {
				continue
			} else if !wentPastStop {
				t := fromMvmt.Timestamp.Timestamp.Sub(reachedStopMvmt.Timestamp.Timestamp.Time) / time.Minute
				fmt.Printf("Time taken: %d\n", t)
			}
		}
	}
}
