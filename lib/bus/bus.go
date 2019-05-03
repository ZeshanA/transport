package bus

import (
	"fmt"
	"log"
	"strings"
	"transport/lib/nulltypes"

	"github.com/lib/pq"

	"gopkg.in/guregu/null.v3"
)

// Constants
const AverageDistanceBetweenStops = 317

// StopDistances
type StopDistance struct {
	RouteID     string
	DirectionID int
	FromID      string
	ToID        string
	Distance    float64
}

func (sd *StopDistance) string() string {
	return fmt.Sprintf("%s – Direction %d – From %s – To %s = %f metres", sd.RouteID, sd.DirectionID, sd.FromID, sd.ToID, sd.Distance)
}

// Internal Format Structs
// ==========================================================================
type VehicleJourney struct {
	LineRef                  null.String
	DirectionRef             null.Int
	TripID                   null.String
	PublishedLineName        null.String
	OperatorRef              null.String
	OriginRef                null.String
	DestinationRef           null.String
	OriginAimedDepartureTime nulltypes.Timestamp
	SituationRef             nulltypes.StringSlice
	Longitude                null.Float
	Latitude                 null.Float
	ProgressRate             null.String
	Occupancy                null.String
	VehicleRef               null.String
	ExpectedArrivalTime      nulltypes.Timestamp
	ExpectedDepartureTime    nulltypes.Timestamp
	DistanceFromStop         null.Int
	NumberOfStopsAway        null.Int
	StopPointRef             null.String
	Timestamp                nulltypes.Timestamp
}

func (vj *VehicleJourney) Value() []interface{} {
	return []interface{}{
		vj.LineRef.String, vj.DirectionRef.Int64, vj.TripID.String, vj.PublishedLineName.String, vj.OperatorRef.String,
		vj.OriginRef.String, vj.DestinationRef.String, vj.OriginAimedDepartureTime,
		pq.Array(vj.SituationRef.StringSlice), vj.Longitude.Float64, vj.Latitude.Float64,
		vj.ProgressRate.String, vj.Occupancy.String, vj.VehicleRef.String, vj.ExpectedArrivalTime,
		vj.ExpectedDepartureTime, vj.DistanceFromStop.Int64, vj.NumberOfStopsAway.Int64,
		vj.StopPointRef.String, vj.Timestamp,
	}
}

type DirectedRoute struct {
	RouteID     string
	DirectionID int
	VehicleRef  string
}

type LabelledJourney struct {
	LineRef               null.String
	DirectionRef          null.Int
	OperatorRef           null.String
	OriginRef             null.String
	DestinationRef        null.String
	Longitude             null.Float
	Latitude              null.Float
	ProgressRate          null.String
	Occupancy             null.String
	VehicleRef            null.String
	ExpectedArrivalTime   nulltypes.Timestamp
	ExpectedDepartureTime nulltypes.Timestamp
	DistanceFromStop      null.Int
	NumberOfStopsAway     null.Int
	StopPointRef          null.String
	Timestamp             nulltypes.Timestamp
	TimeToStop            null.Int
}

func LabelledJourneyFrom(mvmt VehicleJourney, timeToStop int) LabelledJourney {
	return LabelledJourney{
		LineRef:               mvmt.LineRef,
		DirectionRef:          mvmt.DirectionRef,
		OperatorRef:           mvmt.OperatorRef,
		OriginRef:             mvmt.OriginRef,
		DestinationRef:        mvmt.DestinationRef,
		Longitude:             mvmt.Longitude,
		Latitude:              mvmt.Latitude,
		ProgressRate:          mvmt.ProgressRate,
		Occupancy:             mvmt.Occupancy,
		VehicleRef:            mvmt.VehicleRef,
		ExpectedArrivalTime:   mvmt.ExpectedArrivalTime,
		ExpectedDepartureTime: mvmt.ExpectedDepartureTime,
		DistanceFromStop:      mvmt.DistanceFromStop,
		NumberOfStopsAway:     mvmt.NumberOfStopsAway,
		StopPointRef:          mvmt.StopPointRef,
		Timestamp:             mvmt.Timestamp,
		TimeToStop:            null.IntFrom(int64(timeToStop)),
	}
}

// LabelledJourneyToInterface converts a slice of LabelledJourney structs into
// a slice of interface{}
func LabelledJourneyToInterface(journeys []LabelledJourney) []interface{} {
	r := make([]interface{}, len(journeys))
	for i, journey := range journeys {
		r[i] = journey
	}
	return r
}

// ExtractEntriesFromLabelledJourney converts a single LabelledJourney struct into
// a slice of interface{} which represents the database row
func ExtractEntriesFromLabelledJourney(ljEntry interface{}) []interface{} {
	lj, ok := ljEntry.(LabelledJourney)
	if !ok {
		log.Panicf("ExtractEntriesFromLabelledJourney: entry passed in is not an LabelledJourney struct")
	}
	return []interface{}{
		lj.LineRef, lj.DirectionRef, lj.OperatorRef, lj.OriginRef, lj.DestinationRef, lj.Longitude, lj.Latitude,
		lj.ProgressRate, lj.Occupancy, lj.VehicleRef, lj.ExpectedArrivalTime, lj.ExpectedDepartureTime,
		lj.DistanceFromStop, lj.NumberOfStopsAway, lj.StopPointRef, lj.Timestamp, lj.TimeToStop,
	}
}

func PartitionJourneys(journeys []VehicleJourney) map[DirectedRoute][]VehicleJourney {
	result := map[DirectedRoute][]VehicleJourney{}
	for _, journey := range journeys {
		route := DirectedRoute{
			journey.LineRef.String,
			int(journey.DirectionRef.Int64),
			journey.VehicleRef.String,
		}
		result[route] = append(result[route], journey)
	}
	return result
}

func RemoveAgencyID(routeID string) string {
	split := strings.Split(routeID, "_")
	return split[1]
}
