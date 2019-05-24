package request

import (
	"fmt"
	"time"
	"transport/lib/database"
)

type JourneyParams struct {
	RouteID     string
	DirectionID int
	FromStop    string
	ToStop      string
	ArrivalTime time.Time
}

func (jp JourneyParams) MarshalJSON() ([]byte, error) {
	templates := `{
		"routeID": "%s",
		"directionID": "%d",
		"fromStop": "%s",
		"toStop": "%s",
		"arrivalTime": "%s"
	}`
	str := fmt.Sprintf(templates, jp.RouteID, jp.DirectionID, jp.FromStop, jp.ToStop, jp.ArrivalTime.Format(database.TimeFormat))
	return []byte(str), nil
}
