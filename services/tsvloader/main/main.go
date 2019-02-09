package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
)

func main() {
	hostID, err := strconv.Atoi(os.Args[1])
	hostCount, err := strconv.Atoi(os.Args[2])
	if err != nil {
		log.Fatalf("Failed to convert args to integers: %v\n", os.Args)
	}
	fetchAndStoreArchives(hostID, hostCount)
}

// Gets URLs for the mtaArchive date range and concurrently
// fetches and stores the data from each URL
func fetchAndStoreArchives(hostID int, hostCount int) {
	// Get URLs
	URLs := getURLsForDateRange(mtaArchiveStartDate, mtaArchiveEndDate)
	// Number of URLs each host needs to process
	taskCount := len(URLs) / hostCount
	// The index of the first URL this host should process (based on its ID)
	firstTaskIndex := hostID * taskCount

	fmt.Printf("Host ID: %d	Host Count: %d", hostID, hostCount)
	fmt.Printf("First Task Index: %d\n", firstTaskIndex)
	fmt.Printf("Last Task Index: %d\n", firstTaskIndex+taskCount-1)

	// Process 'taskCount' URLs starting from firstTaskIndex
	for i := firstTaskIndex; i < firstTaskIndex+taskCount; i++ {
		fetchAndStore(URLs[i])
	}
}

// Fetches and stores the data for a single URL
func fetchAndStore(URL string) {
	filename := fetchSingleDay(URL)
	decompressedFilename := decompressFile(filename)
	validRows := removeNullRows(decompressedFilename)
	arrivalEntries := unmarshalMTADataBytes(validRows)
	removeDataFiles(filename, decompressedFilename)
	store(arrivalEntries)
}
