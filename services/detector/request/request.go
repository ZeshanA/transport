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

func GetParams() JourneyParams {
	// TODO: These should come from a user request; using constants for now
	routeID, directionID, fromStop, toStop := "MTA NYCT_S78", 1, "MTA_200177", "MTA_201081"
	now := time.Now().In(database.TimeLoc)
	arrivalTime := now.Add(4 * time.Hour)
	return JourneyParams{RouteID: routeID, DirectionID: directionID, FromStop: fromStop, ToStop: toStop, ArrivalTime: arrivalTime}
}
