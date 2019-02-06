package main

import (
	"database/sql/driver"
	"time"
)

// Constants
const timeFormat = "2006-01-02 15:04:05"

var mtaArchiveStartDate = time.Date(2014, 8, 1, 0, 0, 0, 0, time.UTC)
var mtaArchiveEndDate = time.Date(2014, 11, 1, 0, 0, 0, 0, time.UTC)

// ArrivalEntry stores a single parsed bus entry from the MTA data
type ArrivalEntry struct {
	Latitude         float64   `csv:"latitude"`
	Longitude        float64   `csv:"longitude"`
	Timestamp        Timestamp `csv:"time_received"`
	VehicleID        int32     `csv:"vehicle_id"`
	DistanceAlong    float64   `csv:"distance_along_trip"`
	DirectionID      int32     `csv:"inferred_direction_id"`
	Phase            string    `csv:"inferred_phase"`
	RouteID          string    `csv:"inferred_route_id"`
	TripID           string    `csv:"inferred_trip_id"`
	NextStopDistance float64   `csv:"next_scheduled_stop_distance"`
	NextStopID       string    `csv:"next_scheduled_stop_id"`
}

// TableName defines the name of the table in the DB that
// ArrivalEntry structs will be added to
func (ArrivalEntry) TableName() string {
	return "arrival"
}

// Timestamp is a wrapper around time.Time to allow for a custom
// UnmarshalCSV method
type Timestamp struct {
	time.Time
}

// UnmarshalCSV method on ArrivalEntry to specify how
// to unmarshal incoming date strings.
func (t *Timestamp) UnmarshalCSV(csv string) (err error) {
	t.Time, err = time.Parse(timeFormat, csv)
	return err
}

// Value describes how to marshal a Timestamp when inserting into a DB
func (t Timestamp) Value() (driver.Value, error) {
	return t.Time, nil
}

// Scan defines how values fetched from the DB can b
// unmarshalled into a Timestamp struct
// TODO: Implement Scanner interface for timestamps or else
//  structs from DB can't be unmarshalled
func (t *Timestamp) Scan(value interface{}) error {
	return nil
}
