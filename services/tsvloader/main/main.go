package main

import (
	"fmt"
	"net/http"
)

func main() {
	fmt.Println("Starting HTTP server...")
	// When /fetch is loaded, trigger fetching and storing of archives
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
	arrivalEntries := unmarshalMTADataFile(decompressedFilename)

	// Put each struct from unmarshalMTADataFile array into Postgres
	for _, entry := range arrivalEntries {
		// TODO: Replace the printing below with insertion into a DB
		fmt.Printf("%v\n", entry)
	}
}
