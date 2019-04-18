package main

import (
	"fmt"
	"net/url"
	"strings"
	"testing"
	"transport/lib/bustime"
	"transport/lib/testhelper"

	"googlemaps.github.io/maps"

	"github.com/stretchr/testify/assert"
)

func TestGetDistances(t *testing.T) {
	ts := testhelper.ServeMultiResponseMock(mockedDistanceResponses, extractCoordsFromURL)
	defer ts.Close()

	mc, err := maps.NewClient(maps.WithAPIKey("TEST"), maps.WithBaseURL(ts.URL))
	if err != nil {
		assert.Fail(t, fmt.Sprintf("failed to initialise maps client: %s", err))
	}

	expected := stopDistances
	actual := GetDistances(mc, stopDetails)
	assert.Equal(t, expected, actual)
}

func extractCoordsFromURL(fullURL *url.URL) string {
	from := strings.Replace(fullURL.Query().Get("origins"), "00000", "", -1)
	to := strings.Replace(fullURL.Query().Get("destinations"), "00000", "", -1)
	return fmt.Sprintf("%s;%s", from, to)
}

var mockedDistanceResponses = map[string]string{
	"12.3,34.5;67.8,90.1": `{"rows": [{"elements": [{"distance": {"value": 157}}]}], "status": "OK"}`,
	"67.8,90.1;23.4,56.7": `{"rows": [{"elements": [{"distance": {"value": 148}}]}], "status": "OK"}`,
	"14.1,23.5;13.8,83.1": `{"rows": [{"elements": [{"distance": {"value": 127}}]}], "status": "OK"}`,
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
	{routeID: "MTA M1", directionID: 0, fromID: "Stop1", toID: "Stop2", distance: 157},
	{routeID: "MTA M1", directionID: 0, fromID: "Stop2", toID: "Stop3", distance: 148},
	{routeID: "MTA M1", directionID: 1, fromID: "Stop4", toID: "Stop5", distance: 127},
}
