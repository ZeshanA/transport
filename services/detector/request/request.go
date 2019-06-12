package request

import (
	"fmt"
	"strings"
	"time"
	"transport/lib/database"
)

type JourneyParams struct {
	RouteID     string
	DirectionID int
	FromStop    string
	ToStop      string
	ArrivalTime database.Timestamp
	Channel     string
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

func (jp JourneyParams) String() string {
	b, _ := jp.MarshalJSON()
	s := string(b)
	removeTabs := strings.Replace(s, "\t", "", -1)
	removeNewLines := strings.Replace(removeTabs, "\n", " ", -1)
	return string(removeNewLines)
}

func GetParams() JourneyParams {
	// TODO: These should come from a user request; using constants for now
	// routeID, directionID, fromStop, toStop := "MTA NYCT_S78", 1, "MTA_200177", "MTA_201081"
	routeID, directionID, fromStop, toStop := "MTA NYCT_M86+", 0, "MTA_401901", "MTA_401905"
	return JourneyParams{RouteID: routeID, DirectionID: directionID, FromStop: fromStop, ToStop: toStop, ArrivalTime: database.Timestamp{Time: time.Date(2019, 6, 2, 18, 35, 0, 0, database.TimeLoc)}}
}
