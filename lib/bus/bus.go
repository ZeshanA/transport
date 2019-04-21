package bus

import (
	"fmt"
	"transport/lib/nulltypes"

	"github.com/lib/pq"

	"gopkg.in/guregu/null.v3"
)

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
