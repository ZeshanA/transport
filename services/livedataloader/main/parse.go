package main

import (
	"log"
	"net/url"
	"strings"

	"github.com/tidwall/gjson"
)

// The nestParent type represents an item under "MonitoredVehicleJourney" that
// points to another object, instead of a raw value. If one of the fields
// we're filtering by is nested, we need to append the names of its parent
// before looking it up in the MonitoredVehicleJourney. Creating the nestParent type
// avoids us having to use strings for the parent names in the `nestedFields` lookup table.
type nestParent string

const (
	monitoredCall           nestParent = "MonitoredCall"
	framedVehicleJourneyRef nestParent = "FramedVehicleJourneyRef"
	vehicleLocation         nestParent = "VehicleLocation"
)

// nestedFields is a lookup table for fields that we might need to filter by,
// which aren't stored directly under MonitoredVehicleJourney, but are nested
// inside an additional object within MonitoredVehicleJourney. The nestParent value
// is the name of this parent object that we must prepend when looking up the field.
var nestedFields = map[string]nestParent{
	"ArrivalProximityText":   monitoredCall,
	"DistanceFromStop":       monitoredCall,
	"NumberOfStopsAway":      monitoredCall,
	"StopPointRef":           monitoredCall,
	"VisitNumber":            monitoredCall,
	"StopPointName":          monitoredCall,
	"DataFrameRef":           framedVehicleJourneyRef,
	"DatedVehicleJourneyRef": framedVehicleJourneyRef,
	"Longitude":              vehicleLocation,
	"Latitude":               vehicleLocation,
}

// The JSON path in the response where the VehicleActivity array (containing MonitoredVehicleJourney
// objects) can be found
const vehicleActivityPath = "Siri.ServiceDelivery.VehicleMonitoringDelivery.0.VehicleActivity"

// Returns a JSON array (as a string) that contains only the MonitoredVehicleJourney
// items that satisfied the filters passed in (e.g. LineRef="MTA NYCT_B59")
func createVehicleDataResponse(JSONData *string, filters url.Values) *string {

	log.Println("Creating response...")

	// Extract the array of MonitoredVehicleJourney items from the response
	liveVehicleData := gjson.Get(*JSONData, vehicleActivityPath)

	// If there are no filters, return a string of the entire array
	if len(filters) == 0 {
		log.Printf("No filters specified, returning all objects...")
		liveVehicleData := liveVehicleData.String()
		return &liveVehicleData
	}

	log.Println("Filters found...")

	// Return a JSON array (in string format) containing all matching MonitoredVehicleJourneys
	return getJSONArrayOfMatches(&liveVehicleData, filters)
}

// Takes a pointer to a gjson.Result containing an array of MonitoredVehicleJourneys,
// and a map containing filters. Returns a pointer to a JSON array (in string format)
// containing the MonitoredVehicleJourney items that satisfied all the filters.
func getJSONArrayOfMatches(liveVehicleData *gjson.Result, filters url.Values) *string {
	var matches []gjson.Result

	log.Println("Applying filters...")

	liveVehicleData.ForEach(func(key, value gjson.Result) bool {
		if satisfiesFilters(value, filters) {
			matches = append(matches, value)
		}
		return true
	})

	return createJSONArray(&matches)
}

// Takes a pointer to an array of gjson.Result items and returns a pointer to a JSON
// array (in string format) containing said items (e.g. [ "{}", "{}" ] -> "[ {}, {} ]")
func createJSONArray(elements *[]gjson.Result) *string {
	log.Println("Creating JSON array from filtered result...")
	elementStrings := resultsToStrings(elements)
	commaSeparatedStrings := strings.Join(*elementStrings, ",\n")
	JSONArray := "[" + commaSeparatedStrings + "]"
	return &JSONArray
}

// Takes a pointer to an array of gjson.Result items, and returns an array containing
// each Result item converted to a string
func resultsToStrings(elements *[]gjson.Result) *[]string {
	elementStrings := make([]string, len(*elements))
	for i, element := range *elements {
		elementStrings[i] = element.String()
	}
	return &elementStrings
}

// Returns true iff the values of the fields in `item` match the values
// given in the `filters` map
func satisfiesFilters(item gjson.Result, filters url.Values) bool {

	// Loop over each filter
	for filter, expectedVal := range filters {

		// Most fields are nested directly under the MonitoredVehicleJourney object
		prefix := "MonitoredVehicleJourney."

		// Some fields have an extra layer of nesting, so append the additional parent
		// objects to the path we use to look up the field in the JSON data
		if parent, hasParent := nestedFields[filter]; hasParent {
			prefix += string(parent) + "."
		}

		// Look up the field mentioned in the filter (using the right prefix) and verify
		if item.Get(prefix+filter).String() != expectedVal[0] {
			return false
		}

	}
	return true
}
