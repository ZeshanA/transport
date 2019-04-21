package main

import (
	"encoding/json"
	"log"
	"transport/lib/bus"
	"transport/lib/database"
	"transport/lib/nulltypes"

	"gopkg.in/guregu/null.v3"
)

// Takes a JSON string representing an MTAVehicleMonitoringResponse and return
// a slice containing the same data in the internal VehicleJourney format
func convertToIR(jsonString []byte) []bus.VehicleJourney {
	var response MTAVehicleMonitoringResponse
	err := json.Unmarshal(jsonString, &response)
	if err != nil {
		log.Fatalf("error parsing JSON: %v\n", err)
	}

	externalJourneys := response.Siri.ServiceDelivery.VehicleMonitoringDelivery[0].VehicleActivity
	var internalJourneys = make([]bus.VehicleJourney, len(externalJourneys))

	for i, jrny := range externalJourneys {
		mvj, ts := jrny.MonitoredVehicleJourney, jrny.RecordedAtTime
		internalJourneys[i] = getVehicleJourney(mvj, ts)
	}

	return internalJourneys
}

// Converts an MTAMonitoredVehicleJourney into the internal VehicleJourney format
func getVehicleJourney(mvj MTAMonitoredVehicleJourney, timestamp database.Timestamp) bus.VehicleJourney {
	return bus.VehicleJourney{
		LineRef:                  null.StringFrom(mvj.LineRef),
		DirectionRef:             null.IntFrom(int64(mvj.DirectionRef)),
		TripID:                   null.StringFrom(mvj.FramedVehicleJourneyRef.DatedVehicleJourneyRef),
		PublishedLineName:        null.StringFrom(mvj.PublishedLineName[0]),
		OperatorRef:              null.StringFrom(mvj.OperatorRef),
		OriginRef:                null.StringFrom(mvj.OriginRef),
		DestinationRef:           null.StringFrom(mvj.DestinationRef),
		OriginAimedDepartureTime: nulltypes.TimestampFrom(mvj.OriginAimedDepartureTime),
		SituationRef:             nulltypes.StringSliceFrom(flattenSituationRef(mvj.SituationRef)),
		Longitude:                null.FloatFrom(mvj.VehicleLocation.Longitude),
		Latitude:                 null.FloatFrom(mvj.VehicleLocation.Latitude),
		ProgressRate:             null.StringFrom(mvj.ProgressRate),
		Occupancy:                null.StringFrom(mvj.Occupancy),
		VehicleRef:               null.StringFrom(mvj.VehicleRef),
		ExpectedArrivalTime:      nulltypes.TimestampFrom(mvj.MonitoredCall.ExpectedArrivalTime),
		ExpectedDepartureTime:    nulltypes.TimestampFrom(mvj.MonitoredCall.ExpectedDepartureTime),
		DistanceFromStop:         null.IntFrom(int64(mvj.MonitoredCall.DistanceFromStop)),
		NumberOfStopsAway:        null.IntFrom(int64(mvj.MonitoredCall.NumberOfStopsAway)),
		StopPointRef:             null.StringFrom(mvj.MonitoredCall.StopPointRef),
		Timestamp:                nulltypes.TimestampFrom(timestamp),
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
