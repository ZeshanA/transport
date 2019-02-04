package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"time"
	"transport/lib/network"

	"github.com/gocarina/gocsv"
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

// Returns a message telling the sender their request has been received
// and then fetches and stores all archives
func fetchArchivesEndpoint(w http.ResponseWriter, _ *http.Request) {
	fmt.Println("Fetch archives request received...")
	_, err := w.Write([]byte("The server will now fetch data from the MTA, thank you!"))
	if err != nil {
		fmt.Printf("Error in fetchArchivesEndpoint handler: %v", err)
	}
	fetchAndStoreArchives()
}

// Gets URLs for the mtaArchive date range and concurrently
// fetches and stores the data from each URL
func fetchAndStoreArchives() {
	// Get URLs
	URLs := getURLsForDateRange(mtaArchiveStartDate, mtaArchiveEndDate)
	// For each URL
	for _, URL := range URLs {
		go fetchAndStore(URL)
	}
}

// Fetches and stores the data for a single URL
func fetchAndStore(URL string) {
	pathToFile := path.Base(URL)

	// Download the file
	if err := network.DownloadFile(URL, pathToFile); err != nil {
		fmt.Printf("Failed to fetch URL %s due to the following error: %v\n", URL, err)
		return
	}

	// Unzip
	fmt.Println("Succesfully fetched %s", pathToFile)
	// Pass unzipped file path to loadMTAData
	// Shove each struct from loadMTAData array into Postgres
}

// Returns an array of URLs with the date of each day between `start` and `end` (inclusive
// of start and exclusive of end) interpolated into the URL
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

// Takes the MTA data in TSV format and returns an array
// of marshalled ArrivalEntry structs
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
