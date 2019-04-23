package fetch

import (
	"database/sql"
	"fmt"
	"log"
	"time"
	"transport/lib/bus"
	"transport/lib/database"
	"transport/lib/dates"
)

func DateRange(db *sql.DB, startDate time.Time, lastDate time.Time) [][]bus.VehicleJourney {
	endDate := lastDate.AddDate(0, 0, 1)
	var journeys [][]bus.VehicleJourney
	for d := startDate; !dates.Equal(d, endDate); d = d.AddDate(0, 0, 1) {
		data := getDataForDate(db, d)
		journeys = append(journeys, data)
	}
	return journeys
}

func getDataForDate(db *sql.DB, date time.Time) []bus.VehicleJourney {
	rows := queryByDate(date, db)
	defer rows.Close()
	journeys := scanVehicleJournies(rows)
	return journeys
}

func scanVehicleJournies(rows *sql.Rows) []bus.VehicleJourney {
	// Don't need to store the entryID from each row, so just write it
	// to this placeholder and ignore it
	entryIDPtr := 0
	var journeys []bus.VehicleJourney
	for rows.Next() {
		journey := bus.VehicleJourney{}
		err := rows.Scan(
			&journey.LineRef, &journey.DirectionRef, &journey.TripID, &journey.PublishedLineName, &journey.OperatorRef,
			&journey.OriginRef, &journey.DestinationRef, &journey.OriginAimedDepartureTime, &journey.SituationRef,
			&journey.Longitude, &journey.Latitude, &journey.ProgressRate, &journey.Occupancy, &journey.VehicleRef,
			&journey.ExpectedArrivalTime, &journey.ExpectedDepartureTime, &journey.DistanceFromStop,
			&journey.NumberOfStopsAway, &journey.StopPointRef, &journey.Timestamp, &entryIDPtr,
		)
		if err != nil {
			log.Fatalf("getDataForDate: error whilst scanning row from DB into a struct: %s", err)
		}
		journeys = append(journeys, journey)
	}
	err := rows.Err()
	if err != nil {
		log.Fatalf("getDataForDate: error whilst scanning rows from DB: %s\n", err)
	}
	return journeys
}

func queryByDate(date time.Time, db *sql.DB) *sql.Rows {
	q := fmt.Sprintf(
		`SELECT * FROM %[1]s WHERE TIMESTAMP BETWEEN '%[2]d-%[3]d-%[4]d 00:00:00' AND '%[2]d-%[3]d-%[4]d 23:59:59'`,
		database.VehicleJourneyTable.Name,
		date.Year(), date.Month(), date.Day(),
	)
	rows, err := db.Query(q)
	if err != nil {
		log.Fatalf("getDataForDate: error executing SQL query to fetch dates: %s\n", err)
	}
	return rows
}
