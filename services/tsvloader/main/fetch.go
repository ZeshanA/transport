package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"path"
	"path/filepath"
	"time"
	"transport/lib/csvhelper"
	"transport/lib/network"

	"github.com/gocarina/gocsv"
	"github.com/mholt/archiver"
)

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

// Downloads the file from 'URL' and returns the filename it was saved as
func fetchSingleDay(URL string) (filename string) {
	nameOfFile := path.Base(URL)

	// Download the file
	if err := network.DownloadFile(URL, nameOfFile); err != nil {
		panic(fmt.Sprintf("failed to fetch URL %s due to the following error: %v\n", URL, err))
	}

	fmt.Printf("Succesfully fetched %s...\n", nameOfFile)

	return nameOfFile
}

// Decompresses the file at 'path' and stores the result at '${path}_uncompressed'
// e.g. decompressFile("main/abc.xz") => "main/abc.xz_uncompressed"
func decompressFile(path string) (decompressedPath string) {
	newPath := path + "_uncompressed"
	err := archiver.DecompressFile(path, newPath)
	if err != nil {
		panic(fmt.Sprintf("failed to unzip file '%s' due to the following error: %v\n", newPath, err))
	}
	fmt.Printf("Successfully decompressed %s...\n", path)
	return newPath
}

// Removes all rows in the TSV file at 'path' that have null values in any column
// Returns the number of null rows that were removed
func removeNullRows(path string) (nullRowCount int) {
	nullCount, err := csvhelper.RemoveNullRows(path, "\t")
	if err != nil {
		panic(fmt.Sprintf(
			"failed to remove null rows from '%s' due to the following error: %v\n",
			path,
			err,
		))
	}
	fmt.Printf("Successfully removed %d null rows from %s...\n", nullCount, path)
	return nullCount
}

// Takes the MTA data in TSV format and returns an array
// of marshalled ArrivalEntry structs
func unmarshalMTADataFile(path string) []ArrivalEntry {

	// Clean the path passed in
	cleanedPath := filepath.Clean(path)

	fmt.Printf("Loading in rows from %s...\n", cleanedPath)

	// Open up the .tsv file
	arrivalsFile, err := os.OpenFile(cleanedPath, os.O_RDONLY, os.ModePerm)
	if err != nil {
		panic(fmt.Errorf("Error whilst reading MTA data file: \n +%v", err))
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

	return entries
}
