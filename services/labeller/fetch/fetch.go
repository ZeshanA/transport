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
	log.Printf("Fetching rows from vehicle_journey table with timestamps between %s and %s\n", startDate, lastDate)
	endDate := lastDate.AddDate(0, 0, 1)
	var journeys [][]bus.VehicleJourney
	for d := startDate; !dates.Equal(d, endDate); d = d.AddDate(0, 0, 1) {
		data := getDataForDate(db, d)
		journeys = append(journeys, data)
	}
	log.Printf("Succesfully fetched rows for timestamps between %s and %s\n", startDate, lastDate)
	return journeys
}

func getDataForDate(db *sql.DB, date time.Time) []bus.VehicleJourney {
	log.Printf("Fetching rows for date %s\n", date)
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

func queryByDate(start time.Time, db *sql.DB) *sql.Rows {
	end := start.AddDate(0, 0, 1)
	q := fmt.Sprintf(
		`SELECT * FROM %s WHERE TIMESTAMP BETWEEN '%s 04:00:00' AND '%s 03:59:59' ORDER BY TIMESTAMP ASC`,
		database.VehicleJourneyTable.Name,
		start.Format(database.DateFormat),
		end.Format(database.DateFormat),
	)
	rows, err := db.Query(q)
	if err != nil {
		log.Fatalf("getDataForDate: error executing SQL query to fetch dates: %s\n", err)
	}
	return rows
}
