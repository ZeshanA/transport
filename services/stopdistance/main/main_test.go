package main

import (
	"strings"
	"testing"
	"transport/lib/bustime"
	"transport/lib/mapping"
	"transport/lib/testhelper"

	"github.com/stretchr/testify/assert"
)

func TestGetDistances(t *testing.T) {
	ts := testhelper.ServeMultiResponseMock(mockedDistanceResponses, extractCoordsFromURL)
	defer ts.Close()

	expected := stopDistances

	mc := mapping.NewClient("TEST", mapping.CustomBaseURLOption(ts.URL))
	actual := GetDistances(mc, stopDetails)

	assert.Equal(t, expected, actual)
}

func extractCoordsFromURL(url string) string {
	components := strings.Split(url, "/")
	coordinates := components[len(components)-1]
	// Conversion from float64 coordinates to strings adds five 0s onto the end,
	// remove them here just to make hardcoding the mocked responses easier
	return strings.Replace(coordinates, "00000", "", -1)
}

var mockedDistanceResponses = map[string]string{
	"12.3,34.5;67.8,90.1": `{"code": "Ok", "distances": [[157.8]]}`,
	"67.8,90.1;23.4,56.7": `{"code": "Ok", "distances": [[148.1]]}`,
	"14.1,23.5;13.8,83.1": `{"code": "Ok", "distances": [[127.2]]}`,
}

var stopDetails = map[string]map[int][]bustime.BusStop{
	"MTA M1": {
		0: {
			bustime.BusStop{ID: "Stop1", Latitude: 12.3, Longitude: 34.5},
			bustime.BusStop{ID: "Stop2", Latitude: 67.8, Longitude: 90.1},
			bustime.BusStop{ID: "Stop3", Latitude: 23.4, Longitude: 56.7},
		},
		1: {
			bustime.BusStop{ID: "Stop4", Latitude: 14.1, Longitude: 23.5},
			bustime.BusStop{ID: "Stop5", Latitude: 13.8, Longitude: 83.1},
		},
	},
	"MTA M2": {
		0: {
			bustime.BusStop{ID: "Stop6", Latitude: 10.3, Longitude: 13.5},
		},
		1: {},
	},
}

var stopDistances = []stopDistance{
	{routeID: "MTA M1", directionID: 0, fromID: "Stop1", toID: "Stop2", distance: 157.8},
	{routeID: "MTA M1", directionID: 0, fromID: "Stop2", toID: "Stop3", distance: 148.1},
	{routeID: "MTA M1", directionID: 1, fromID: "Stop4", toID: "Stop5", distance: 127.2},
}
