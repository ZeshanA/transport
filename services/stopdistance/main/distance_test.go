package main

import (
	"fmt"
	"log"
	"net/url"
	"sort"
	"strings"
	"testing"
	"transport/lib/bus"
	"transport/lib/bustime"
	"transport/lib/testhelper"

	"github.com/stretchr/testify/assert"

	"googlemaps.github.io/maps"
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
	sort.Slice(actual, func(i, j int) bool {
		return actual[i].FromID < actual[j].FromID
	})
	assert.Equal(t, expected, actual)
}

func extractCoordsFromURL(fullURL *url.URL) string {
	from, err := url.PathUnescape(strings.Replace(fullURL.Query().Get("origins"), "00000", "", -1))
	if err != nil {
		log.Panicf("extractCoordsFromURL: %s", err)
	}
	to, err := url.PathUnescape(strings.Replace(fullURL.Query().Get("destinations"), "00000", "", -1))
	if err != nil {
		log.Panicf("extractCoordsFromURL: %s", err)
	}
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

var stopDistances = []bus.StopDistance{
	{RouteID: "MTA M1", DirectionID: 0, FromID: "Stop1", ToID: "Stop2", Distance: 157},
	{RouteID: "MTA M1", DirectionID: 0, FromID: "Stop2", ToID: "Stop3", Distance: 148},
	{RouteID: "MTA M1", DirectionID: 1, FromID: "Stop4", ToID: "Stop5", Distance: 127},
}
