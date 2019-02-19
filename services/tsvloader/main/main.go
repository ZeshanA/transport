package main

import (
	"log"
	"os"
	"strconv"
)

func main() {
	hostID, err := strconv.Atoi(os.Args[1])
	hostCount, err := strconv.Atoi(os.Args[2])
	storageDirectory := os.Args[3]
	if err != nil {
		log.Fatalf("Failed to convert args to integers: %v\n", os.Args)
	}
	err = os.MkdirAll(storageDirectory, os.ModePerm)
	if err != nil {
		log.Fatalf("Failed to create storage directory %s due to: %v", storageDirectory, err)
	}
	fetchAndStoreArchives(hostID, hostCount, storageDirectory)
}

// Gets URLs for the mtaArchive date range and concurrently
// fetches and stores the data from each URL
func fetchAndStoreArchives(hostID int, hostCount int, storageDirectory string) {
	// Get URLs
	URLs := getURLsForDateRange(mtaArchiveStartDate, mtaArchiveEndDate)
	// Number of URLs each host needs to process
	taskCount := len(URLs) / hostCount
	// The index of the first URL this host should process (based on its ID)
	firstTaskIndex := hostID * taskCount

	log.Printf("Host ID: %d	Host Count: %d\n", hostID, hostCount)
	log.Printf("First Task Index: %d\n", firstTaskIndex)
	log.Printf("Last Task Index: %d\n", firstTaskIndex+taskCount-1)

	// Process 'taskCount' URLs starting from firstTaskIndex
	for i := firstTaskIndex; i < firstTaskIndex+taskCount; i++ {
		fetchAndStore(URLs[i], storageDirectory)
	}
}

// Fetches and stores the data for a single URL
func fetchAndStore(URL string, storageDirectory string) {
	compressedFile := fetchSingleDay(URL, storageDirectory)
	decompressedFile := decompressFile(compressedFile)
	validRows := removeNullRows(decompressedFile)
	arrivalEntries := unmarshalMTADataBytes(validRows)
	removeDataFiles(compressedFile, decompressedFile)
	store(arrivalEntries)
}
