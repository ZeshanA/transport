package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/gocarina/gocsv"
)

const historicalDataPath = "../data/sample.tsv"
const timeFormat = "2006-01-02 15:04:05"

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

// Timestamp is a wrapper around time.Time to allow for a custom UnmarshalCSV method
type Timestamp time.Time

// UnmarshalCSV method on ArrivalEntry to specify how
// to unmarshal incoming date strings.
func (t *Timestamp) UnmarshalCSV(csv string) (err error) {
	parsedTime, err := time.Parse(timeFormat, csv)
	*t = Timestamp(parsedTime)
	return err
}

func loadMTAData(path string) ([]ArrivalEntry, error) {
	cleanedPath := filepath.Clean(path)
	arrivalsFile, err := os.OpenFile(cleanedPath, os.O_RDONLY, os.ModePerm)
	if err != nil {
		return nil, fmt.Errorf("Error whilst reading MTA data file: \n +%v", err)
	}
	defer arrivalsFile.Close()

	var entries []ArrivalEntry

	gocsv.SetCSVReader(func(in io.Reader) gocsv.CSVReader {
		r := csv.NewReader(in)
		r.Comma = '\t'
		return r
	})

	if err := gocsv.UnmarshalFile(arrivalsFile, &entries); err != nil {
		panic(err)
	}

	return entries, nil
}

func sayHello(w http.ResponseWriter, r *http.Request) {
	message := r.URL.Path
	message = strings.TrimPrefix(message, "/")
	message = "Hello " + message
	w.Write([]byte(message))
}
func main() {
	http.HandleFunc("/", sayHello)
	if err := http.ListenAndServe(":80", nil); err != nil {
		panic(err)
	}
}
