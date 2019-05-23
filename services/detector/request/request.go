package request

import "time"

type JourneyParams struct {
	RouteID     string
	DirectionID int
	FromStop    string
	ToStop      string
	ArrivalTime time.Time
}
