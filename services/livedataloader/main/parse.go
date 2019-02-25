package main

import (
	"encoding/json"
	"log"
	"net/url"
	"strings"
	"time"

	"github.com/tidwall/gjson"
)

// Takes a JSON string representing an MTAVehicleMonitoringResponse and return
// a JSON string with the same data in the internal VehicleJourney format
func convertToIR(jsonString []byte) string {
	var response MTAVehicleMonitoringResponse
	err := json.Unmarshal(jsonString, &response)
	if err != nil {
		log.Fatalf("error parsing JSON: %v\n", err)
	}
	result, err := json.Marshal(response)
	if err != nil {
		log.Fatalf("error marshalling JSON for response: %v\n", err)
	}
	return string(result)
}

// MarshalJSON is a custom marshalling method for MTAVehicleMonitoringResponse
// that converts the response to the internal data format and performs marshalling
// on that struct
func (response MTAVehicleMonitoringResponse) MarshalJSON() ([]byte, error) {
	vehicleActivity := response.Siri.ServiceDelivery.VehicleMonitoringDelivery[0].VehicleActivity
	var journeys = make([]VehicleJourney, len(vehicleActivity))
	for i, activityEntry := range vehicleActivity {
		MTAJourney := activityEntry.MonitoredVehicleJourney
		timestamp := activityEntry.RecordedAtTime
		journeys[i] = getVehicleJourney(MTAJourney, timestamp)
	}
	return json.Marshal(journeys)
}

// Converts an MTAMonitoredVehicleJourney into the internal VehicleJourney format
func getVehicleJourney(mvj MTAMonitoredVehicleJourney, timestamp Timestamp) VehicleJourney {
	return VehicleJourney{
		mvj.LineRef,
		mvj.DirectionRef,
		mvj.FramedVehicleJourneyRef.DatedVehicleJourneyRef,
		mvj.PublishedLineName[0],
		mvj.OperatorRef,
		mvj.OriginRef,
		mvj.DestinationRef,
		mvj.OriginAimedDepartureTime,
		flattenSituationRef(mvj.SituationRef),
		mvj.VehicleLocation.Longitude,
		mvj.VehicleLocation.Latitude,
		mvj.ProgressRate,
		mvj.Occupancy,
		mvj.VehicleRef,
		mvj.MonitoredCall.ExpectedArrivalTime,
		mvj.MonitoredCall.ExpectedDepartureTime,
		mvj.MonitoredCall.DistanceFromStop,
		mvj.MonitoredCall.NumberOfStopsAway,
		mvj.MonitoredCall.StopPointRef,
		timestamp,
	}
}

// Converts a slice of MTASituationRef into a slice of strings
// representing just the IDs found in MTASituationRef
func flattenSituationRef(refs []MTASituationRef) []string {
	var flattened = make([]string, len(refs))
	for i, ref := range refs {
		flattened[i] = ref.SituationSimpleRef
	}
	return flattened
}

// Timestamp is a wrapper around time.Time to allow for a custom
// UnmarshalJSON method
type Timestamp struct {
	time.Time
}

// Custom parsing of incoming timestamps
func (t *Timestamp) UnmarshalJSON(b []byte) error {
	noQuotes := strings.Replace(string(b), "\"", "", 2)
	parsed, err := time.Parse(time.RFC3339, noQuotes)
	if err != nil {
		log.Printf("error whilst parsing Timestamp: %v", err)
	}
	*t = Timestamp{parsed}
	return nil
}

/***************************************************************************************************/

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

// VehicleActivityPath is the JSON path in the response where the VehicleActivity array
// (containing MonitoredVehicleJourney objects) can be found
const VehicleActivityPath = "Siri.ServiceDelivery.VehicleMonitoringDelivery.0.VehicleActivity"

// Returns a JSON array (as a string) that contains only the MonitoredVehicleJourney
// items that satisfied the filters passed in (e.g. LineRef="MTA NYCT_B59")
func createVehicleDataResponse(JSONData *string, filters url.Values) *string {

	log.Println("Creating response...")

	// Extract the array of MonitoredVehicleJourney items from the response
	liveVehicleData := gjson.Get(*JSONData, VehicleActivityPath)

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

	liveVehicleData.ForEach(func(_, value gjson.Result) bool {
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
