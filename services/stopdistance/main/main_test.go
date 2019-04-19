package main

import (
	"fmt"
	"log"
	"net/url"
	"sort"
	"strings"
	"testing"
	"transport/lib/bustime"
	"transport/lib/testhelper"

	"github.com/stretchr/testify/assert"
	"googlemaps.github.io/maps"
)

func TestGetDistancesIntegration(t *testing.T) {
	btMock := testhelper.ServeMultiResponseMock(bustimeResponses, extractBustimeEndpoint)
	bt := bustime.NewClient("TEST", bustime.CustomBaseURLOption(btMock.URL))

	mapsMock := testhelper.ServeMultiResponseMock(distanceResponses, extractCoordsFromURL)
	mc, err := maps.NewClient(maps.WithAPIKey("TEST"), maps.WithBaseURL(mapsMock.URL))
	if err != nil {
		log.Panicf("main: failed to initialise Maps API client: %s", err)
	}

	agencies := bt.GetAgencies()
	log.Printf("%d agencies fetched\n", len(agencies))
	routes := bt.GetRoutes(agencies...)
	log.Printf("%d routes fetched\n", len(agencies))
	stopDetails := bt.GetStops(routes...)

	// Calculate distances between stops and store in DB
	expected := []stopDistance{
		{routeID: "MTA NYCT_M1", directionID: 0, fromID: "MTA_100001", toID: "MTA_100002", distance: 123},
		{routeID: "MTA NYCT_M1", directionID: 0, fromID: "MTA_100002", toID: "MTA_100003", distance: 456},
		{routeID: "MTA NYCT_M1", directionID: 1, fromID: "MTA_100004", toID: "MTA_100005", distance: 789},
		{routeID: "MTA NYCT_M1", directionID: 1, fromID: "MTA_100005", toID: "MTA_100006", distance: 1011},
	}
	actual := GetDistances(mc, stopDetails)
	sort.Slice(actual, func(i, j int) bool {
		return actual[i].fromID < actual[j].fromID
	})
	assert.Equal(t, expected, actual)
}

func extractBustimeEndpoint(fullURL *url.URL) string {
	if strings.Contains(fullURL.String(), "agencies-with-coverage") {
		return "agencies"
	} else if strings.Contains(fullURL.String(), "routes-for-agency") {
		return "routes"
	} else if strings.Contains(fullURL.String(), "stops-for-route") {
		return "stops"
	} else {
		return "invalid"
	}
}

func distanceJSON(distance int) string {
	return fmt.Sprintf(`{"rows": [{"elements": [{"distance": {"value": %d}}]}], "status": "OK"}`, distance)
}

var distanceResponses = map[string]string{
	"40.7,-73.9;40.8,-73.3": distanceJSON(123),
	"40.8,-73.3;40.1,-72.8": distanceJSON(456),
	"49.7,-73.9;48.8,-72.3": distanceJSON(789),
	"48.8,-72.3;49.1,-73.8": distanceJSON(1011),
}

var bustimeResponses = map[string]string{
	"agencies": `{"data": {"list":[{"agencyId": "MTA NYCT"}]}}`,
	"routes":   `{"data": {"list": [{"id": "MTA NYCT_M1"}]}}`,
	"stops":    stopsResponse,
}

var stopsResponse = `
{
  "data": {
    "entry": {
      "stopGroupings": [
        {
          "stopGroups": [
            {
              "id": 0,
              "stopIds": [
                "MTA_100001",
                "MTA_100002",
                "MTA_100003"
              ]
            },
            {
              "id": 1,
              "stopIds": [
                "MTA_100004",
                "MTA_100005",
                "MTA_100006"
              ]
            }
          ]
        }
      ]
    },
    "references": {
      "stops": [
        {
          "id": "MTA_100001",
          "lat": "40.7",
          "lon": "-73.9"
        },
        {
          "id": "MTA_100002",
          "lat": "40.8",
          "lon": "-73.3"
        },
        {
          "id": "MTA_100003",
          "lat": "40.1",
          "lon": "-72.8"
        },
        {
          "id": "MTA_100004",
          "lat": "49.7",
          "lon": "-73.9"
        },
        {
          "id": "MTA_100005",
          "lat": "48.8",
          "lon": "-72.3"
        },
        {
          "id": "MTA_100006",
          "lat": "49.1",
          "lon": "-73.8"
        }
      ]
    }
  }
}`
