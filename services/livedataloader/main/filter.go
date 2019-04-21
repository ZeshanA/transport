package main

import (
	"encoding/json"
	"log"
	"net/url"
	"strings"
	"transport/lib/bus"

	"github.com/tidwall/gjson"
)

// Returns a JSON array (as a string) that contains only the VehicleJourney
// items that satisfied the filters passed in (e.g. LineRef="MTA NYCT_B59")
func createVehicleDataResponse(journeys []bus.VehicleJourney, filters url.Values) []byte {
	log.Println("Creating response...")

	// If there are no filters, return the original data
	if len(filters) == 0 {
		log.Printf("No filters specified, returning all objects...")
		resp, err := json.Marshal(journeys)
		if err != nil {
			log.Printf("createVehicleDataResponse: error marshalling journeys into JSON: %s\n", err)
			return nil
		}
		return resp
	}

	log.Printf("Filters requested: %s\n", filters)
	// Return a JSON array (in string format) containing all matching MonitoredVehicleJourneys
	return getJSONArrayOfMatches(journeys, filters)
}

// Takes a pointer to a gjson.Result containing an array of VehicleJourneys,
// and a map containing filters. Returns a pointer to a JSON array (in string format)
// containing the VehicleJourney items that satisfied all the filters.
func getJSONArrayOfMatches(liveVehicleData []bus.VehicleJourney, filters url.Values) []byte {
	journeysJSON, err := json.Marshal(liveVehicleData)
	if err != nil {
		log.Printf("getJSONArrayOfMatches: error marshalling journeys into JSON: %s\n", err)
		return nil
	}
	parsedJSON := gjson.ParseBytes(journeysJSON)

	var matches []gjson.Result

	log.Println("Applying filters...")

	parsedJSON.ForEach(func(_, value gjson.Result) bool {
		if satisfiesFilters(value, filters) {
			matches = append(matches, value)
		}
		return true
	})

	return createJSONArray(matches)
}

// Takes a pointer to an array of gjson.Result items and returns a JSON array (in string format)
// containing said items (e.g. [ "{}", "{}" ] -> "[ {}, {} ]")
func createJSONArray(elements []gjson.Result) []byte {
	log.Println("Creating JSON array from filtered result...")
	elementStrings := resultsToStrings(elements)
	commaSeparatedStrings := strings.Join(elementStrings, ",\n")
	JSONArray := "[" + commaSeparatedStrings + "]"
	return []byte(JSONArray)
}

// Takes a pointer to an array of gjson.Result items, and returns an array containing
// each Result item converted to a string
func resultsToStrings(elements []gjson.Result) []string {
	elementStrings := make([]string, len(elements))
	for i, element := range elements {
		elementStrings[i] = element.String()
	}
	return elementStrings
}

// Returns true iff the values of the fields in `item` match the values
// given in the `filters` map
func satisfiesFilters(item gjson.Result, filters url.Values) bool {
	for filter, expectedVal := range filters {
		// Look up the field mentioned in the filter (using the right prefix) and return
		// false if it doesn't have the expected value
		if item.Get(filter).String() != expectedVal[0] {
			return false
		}
	}
	return true
}
