package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strings"
	"time"
	"transport/lib/csvhelper"
	"transport/lib/network"

	"github.com/gocarina/gocsv"
	"github.com/mholt/archiver"
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
	filename := fetchSingleDay(URL)
	decompressedFilename := decompressFile(filename)
	removeNullRows(decompressedFilename)
	arrivalEntries := unmarshalFile(decompressedFilename)

	// Put each struct from loadMTAData array into Postgres
	for _, entry := range arrivalEntries {
		// TODO: Replace the printing below with insertion into a DB
		fmt.Printf("%v\n", entry)
	}
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

func fetchSingleDay(URL string) (filename string) {
	nameOfFile := path.Base(URL)

	// Download the file
	if err := network.DownloadFile(URL, nameOfFile); err != nil {
		panic(fmt.Sprintf("failed to fetch URL %s due to the following error: %v\n", URL, err))
	}

	fmt.Printf("Succesfully fetched %s...\n", nameOfFile)

	return nameOfFile
}

func decompressFile(filename string) (decompressedFilename string) {
	// Decompress the .xz file into "filename.xz_uncompressed"
	newFilename := filename + "_uncompressed"
	err := archiver.DecompressFile(filename, newFilename)
	if err != nil {
		panic(fmt.Sprintf("failed to unzip file '%s' due to the following error: %v\n", newFilename, err))
	}

	fmt.Printf("Successfully decompressed %s...\n", filename)

	return newFilename
}

func removeNullRows(filename string) (nullRowCount int) {
	nullCount, err := csvhelper.RemoveNullRows(filename, "\t")
	if err != nil {
		panic(fmt.Sprintf(
			"failed to remove null rows from '%s' due to the following error: %v\n",
			filename,
			err,
		))
	}
	fmt.Printf("Successfully removed %d null rows from %s...\n", nullCount, filename)
	return nullCount
}

func unmarshalFile(filename string) []ArrivalEntry {
	// Pass decompressed file path to loadMTAData
	arrivalEntries, err := loadMTAData(filename)
	if err != nil {
		filenameWithoutExtension := strings.TrimSuffix(filename, filepath.Ext(filename))
		panic(fmt.Sprintf("failed to parse ArrivalEntry structs from %s due to: %v\n", filenameWithoutExtension, err))
	}

	fmt.Printf("Succesfully unmarshalled %s into %d ArrivalEntry structs...\n", filename, len(arrivalEntries))
	return arrivalEntries
}

// Takes the MTA data in TSV format and returns an array
// of marshalled ArrivalEntry structs
func loadMTAData(path string) ([]ArrivalEntry, error) {

	// Clean the path passed in
	cleanedPath := filepath.Clean(path)

	fmt.Printf("Loading in rows from %s...\n", cleanedPath)

	// Open up the .tsv file
	arrivalsFile, err := os.OpenFile(cleanedPath, os.O_RDONLY, os.ModePerm)
	if err != nil {
		return nil, fmt.Errorf("Error whilst reading MTA data file: \n +%v", err)
	}
	defer arrivalsFile.Close()

	// Our structs will be added to the 'entries' slice as the file is being unmarshalled
	var entries []ArrivalEntry

	// Tell gocsv we're using tabs (TSVs) instead of commas (CSVs)
	gocsv.SetCSVReader(func(in io.Reader) gocsv.CSVReader {
		r := csv.NewReader(in)
		r.Comma = '\t'
		return r
	})

	fmt.Printf("Unmarshalling rows from %s...\n", cleanedPath)

	// Unmarshal the .tsv file into an array of ArrivalEntry structs
	if err := gocsv.UnmarshalFile(arrivalsFile, &entries); err != nil {
		panic(err)
	}

	fmt.Printf("Succesfully unmarshalled %d rows from %s...\n", len(entries), cleanedPath)

	return entries, nil
}
