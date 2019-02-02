package main

import (
	"reflect"
	"testing"
	"time"
)

// expectedStructs are what the ArrivalEntry structs from "sampledata_test.tsv"
// should look like after parsing
var expectedStructs = []ArrivalEntry{{
	Latitude:         40.820351,
	Longitude:        -73.851514,
	Timestamp:        Timestamp(time.Date(2014, time.August, 2, 10, 30, 0, 0, time.UTC)),
	VehicleID:        275,
	DistanceAlong:    12337.92990815825,
	DirectionID:      0,
	Phase:            "LAYOVER_DURING",
	RouteID:          "MTA NYCT_BX36",
	TripID:           "MTA NYCT_WF_C4-Saturday-033000_BX36_102",
	NextStopDistance: 0,
	NextStopID:       "MTA_102978",
}, {
	Latitude:         40.612174,
	Longitude:        -74.035670,
	Timestamp:        Timestamp(time.Date(2014, time.August, 2, 10, 34, 11, 0, time.UTC)),
	VehicleID:        453,
	DistanceAlong:    56.41863532032585,
	DirectionID:      0,
	Phase:            "LAYOVER_DURING",
	RouteID:          "MTA NYCT_B37",
	TripID:           "MTA NYCT_JG_C4-Saturday-039500_B35_7",
	NextStopDistance: 346.4537918200367,
	NextStopID:       "MTA_301722",
}}

func TestLoadMTAData(t *testing.T) {
	parsedData, err := loadMTAData("./sampledata_test.tsv")

	if err != nil {
		t.Errorf("Error whilst loading MTA data: %+v", err)
	}

	for i, expectedEntry := range expectedStructs {
		if !reflect.DeepEqual(parsedData[i], expectedEntry) {
			t.Errorf(
				"Mismatch between entries:\n Expected: %+v\n Received: %+v",
				expectedEntry,
				parsedData[i],
			)
		}
	}
}
