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
	// routeID, directionID, fromStop, toStop := "MTA NYCT_S78", 1, "MTA_200177", "MTA_201081"
	routeID, directionID, fromStop, toStop := "MTA NYCT_M86+", 0, "MTA_401901", "MTA_404681"
	now := time.Now().In(database.TimeLoc)
	arrivalTime := now.Add(15 * time.Minute)
	return JourneyParams{RouteID: routeID, DirectionID: directionID, FromStop: fromStop, ToStop: toStop, ArrivalTime: arrivalTime}
}
