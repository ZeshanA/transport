package fetch

import (
	"database/sql"
	"detector/request"
	"fmt"
	"log"
	"time"
	"transport/lib/bus"
	"transport/lib/bustime"

	"github.com/lib/pq"
)

const arrivalWindow = 2 * time.Hour

func MovementsInWindow(db *sql.DB, stopList []bustime.BusStop, jp request.JourneyParams) ([]bus.LabelledJourney, error) {
	fromHour, toHour := jp.ArrivalTime.Add(-arrivalWindow).Hour(), jp.ArrivalTime.Add(arrivalWindow).Hour()
	log.Printf("Fetching movements in window: %d to %d", fromHour, toHour)
	// TODO: Speed up this query
	query := fmt.Sprintf(`
		SELECT * FROM labelled_journey
		WHERE line_ref='%s' AND direction_ref='%d' AND
		(EXTRACT(hour FROM timestamp) BETWEEN '%d' AND 24 OR EXTRACT(hour FROM timestamp) BETWEEN 0 AND '%d')
		AND stop_point_ref=ANY($1) ORDER BY timestamp ASC;
	`, jp.RouteID, jp.DirectionID, fromHour, toHour)
	// Remove stops on the route that are before the 'fromStop'
	trimmedStopList := bustime.ExtractStops("after", jp.FromStop, true, stopList)
	// Execute query, passing in the trimmedStopList as a Postgres array
	rows, err := db.Query(query, pq.Array(trimmedStopList))
	if err != nil {
		log.Printf("fetch.GetMovementsInWindow: error executing query: %v", err)
		return nil, err
	}
	defer rows.Close()
	// Scan Db rows into journey structs
	journeys, err := bus.ScanLabelledJourneyRows(rows)
	if err != nil {
		return nil, err
	}
	log.Printf("Succesfully fetched movements")
	return journeys, nil
}
