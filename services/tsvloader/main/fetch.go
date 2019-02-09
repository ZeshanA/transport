package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"os"
	"path"
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
	fmt.Printf("Decompressing %s into %s...\n", path, newPath)
	err := archiver.DecompressFile(path, newPath)
	if err != nil {
		panic(fmt.Sprintf("failed to unzip file '%s' due to the following error: %v\n", newPath, err))
	}
	fmt.Printf("Successfully decompressed %s...\n", path)
	return newPath
}

// removeNullRows removes all rows in the TSV file at 'path' that have null values in any column.
// Returns a pointer to a []byte containing only the valid rows.
func removeNullRows(path string) (validRows *[]byte) {
	cleanedRows, err := csvhelper.RemoveNullRows(path, "\t")
	if err != nil {
		panic(fmt.Sprintf(
			"failed to remove null rows from '%s' due to the following error: %v\n",
			path,
			err,
		))
	}
	fmt.Printf("Successfully removed null rows from %s...\n", path)
	return cleanedRows
}

// Takes the MTA data in TSV format (as a []byte) and returns an array
// of marshalled ArrivalEntry structs
func unmarshalMTADataBytes(bytes *[]byte) []ArrivalEntry {
	// Our structs will be added to the 'entries' slice as the bytes are being unmarshalled
	var entries []ArrivalEntry

	// Tell gocsv we're using tabs (TSVs) instead of commas (CSVs)
	gocsv.SetCSVReader(func(in io.Reader) gocsv.CSVReader {
		r := csv.NewReader(in)
		r.Comma = '\t'
		return r
	})

	fmt.Println("Unmarshalling rows...")

	// Unmarshal the .tsv file into an array of ArrivalEntry structs
	if err := gocsv.UnmarshalBytes(*bytes, &entries); err != nil {
		fmt.Printf("Error occurred whilst unmarshalling the following: %s\n", string(*bytes))
		panic(err)
	}

	fmt.Printf("Succesfully unmarshalled %d rows...\n", len(entries))

	return entries
}

// removeDataFiles deletes the files at every path passed in.
// Not recursive (only works for files or *empty* directories).
// Fatal error if any of the deletions fail.
func removeDataFiles(paths ...string) {
	for _, filename := range paths {
		err := os.Remove(filename)
		if err != nil {
			log.Fatalf("failed to delete file %s due to: %v\n", err)
		}
	}
}
