package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/gocarina/gocsv"
)

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

func main() {
	fmt.Println("Starting HTTP server...")
	http.HandleFunc("/fetch", fetchArchivesEndpoint)
	if err := http.ListenAndServe(":80", nil); err != nil {
		panic(err)
	}
}

func fetchArchivesEndpoint(w http.ResponseWriter, _ *http.Request) {
	fmt.Printf("Fetch archives request received...")
	_, err := w.Write([]byte("The server will now fetch data from the MTA, thank you!"))
	if err != nil {
		fmt.Printf("Error in fetchArchivesEndpoint handler: %v", err)
	}
	fetchAndStoreArchives()
}

func fetchAndStoreArchives() {
	fmt.Println("Fetching archives!")
}

func getURLsForDateRange(start, end time.Time) (urls []string) {
	const URLFormat = "http://s3.amazonaws.com/MTABusTime/AppQuest3/MTA-Bus-Time_.%s.txt.xz"
	var URLs []string

	for d := start; !end.Equal(d); d = d.AddDate(0, 0, 1) {
		if d.IsZero() {
			break
		}
		URLs = append(URLs, fmt.Sprintf(URLFormat, d.Format("2006-01-02")))
	}
	return URLs
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
