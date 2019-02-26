package main

import (
	"encoding/json"
	"log"
	"strings"
	"time"
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
