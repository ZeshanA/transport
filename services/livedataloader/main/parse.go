package main

import (
	"fmt"
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

func createVehicleDataResponse(JSONData string, filters url.Values) string {
	result := gjson.Get(JSONData, vehicleActivityPath)
	if len(filters) == 0 {
		return result.String()
	}
	var matching []gjson.Result
	result.ForEach(func(key, value gjson.Result) bool {
		if satisfiesFilters(value, filters) {
			fmt.Printf("Satisfied!")
			matching = append(matching, value)
		}
		return true
	})
	return createJSONArray(matching)
}

func createJSONArray(elements []gjson.Result) string {
	elementStrings := make([]string, len(elements))
	for i, element := range elements {
		elementStrings[i] = element.String()
	}
	commaSeparatedStrings := strings.Join(elementStrings, ",\n")
	return "[" + commaSeparatedStrings + "]"
}

func satisfiesFilters(value gjson.Result, filters url.Values) bool {

	for filter, expectedVal := range filters {
		prefix := "MonitoredVehicleJourney."
		if parent, hasParent := nestedFields[filter]; hasParent {
			prefix += string(parent) + "."
		}

		if value.Get(prefix+filter).String() != expectedVal[0] {
			return false
		}
	}
	return true
}
