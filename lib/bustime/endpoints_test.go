package bustime_test

import (
	"fmt"
	"reflect"
	"testing"
	"transport/lib/bustime"
	"transport/lib/testhelper"
)

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
// Agencies
////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
func TestClient_GetAgencies(t *testing.T) {
	// Create HTTP server to serve mock JSON response
	agencyList := `[{"agencyId": "MTA NYCT"}, {"agencyId": "MTABC"}]`
	response := fmt.Sprintf(`{"data": {"list": %s}}`, agencyList)
	ts := testhelper.ServeMock(response)
	defer ts.Close()

	// Create bustime.client
	client := bustime.NewClient("TEST", bustime.CustomBaseURLOption(ts.URL))

	// We expect the following slice of agency ID strings
	expected := []string{"MTA NYCT", "MTABC"}
	actual := client.GetAgencies()
	if !reflect.DeepEqual(expected, actual) {
		t.Errorf("bustime.GetAgencies did not return expected list of agencyIDs (expected: %s, received: %s)", expected, actual)
	}
}

func TestClient_GetAgenciesEmptyResponse(t *testing.T) {
	// Create HTTP server to serve mock JSON response
	response := fmt.Sprintf(`{"data": {"list": %s}}`, "[]")
	ts := testhelper.ServeMock(response)
	defer ts.Close()

	// Create bustime.client
	client := bustime.NewClient("TEST", bustime.CustomBaseURLOption(ts.URL))

	// We expect nil (no slice should have been allocated, as there were no agency IDs returned)
	var expected []string
	actual := client.GetAgencies()
	if !reflect.DeepEqual(expected, actual) {
		t.Errorf("bustime.GetAgencies did not return empty list of agencyIDs (expected: %s, received: %s)", expected, actual)
	}
}

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
// Routes
////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
func TestClient_GetRoutesWithURL(t *testing.T) {
	// The agency IDs that we want the routes for
	agencyIDsList := []string{"MTA NYCT", "MTABC"}

	// Create HTTP server to serve the correct JSON response, based on the routeID in the request URL
	responses := map[string]string{
		"MTA NYCT": `{"data": {"list": [{"id": "MTA NYCT_BX4A"}, {"id": "MTA NYCT_M31"}]}}`,
		"MTABC":    `{"data": {"list": [{"id": "MTABC_BX4A"}, {"id": "MTABC_M31"}]}}`,
	}
	ts := testhelper.ServeMultiResponseMock(responses)
	defer ts.Close()

	// Create bustime.client
	client := bustime.NewClient("TEST", bustime.CustomBaseURLOption(ts.URL))

	expected := []string{"MTA NYCT_BX4A", "MTA NYCT_M31", "MTABC_BX4A", "MTABC_M31"}
	actual := client.GetRoutes(agencyIDsList...)

	if !reflect.DeepEqual(expected, actual) {
		t.Errorf("bustime.GetRoutes did not return expected list of routeIDs (expected: %s, received: %s)", expected, actual)
	}
}

func TestClient_GetRoutesWithURLEmpty(t *testing.T) {
	// Construct mocked JSON response with empty list
	response := `{"data": {"list": []}}`
	ts := testhelper.ServeMock(response)
	defer ts.Close()

	// Create bustime.client
	client := bustime.NewClient("TEST", bustime.CustomBaseURLOption(ts.URL))

	var expected []string
	actual := client.GetRoutes()
	if !reflect.DeepEqual(expected, actual) {
		t.Errorf("bustime.GetRoutes did not return empty list of routeIDs (expected: %s, received: %s)", expected, actual)
	}
}

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
// Stops
////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

func TestClient_GetStops(t *testing.T) {
	// Create HTTP server to serve the correct JSON response, based on the routeID in the request URL
	ts := testhelper.ServeMultiResponseMock(stopsResponses)
	defer ts.Close()

	// Create bustime.client
	client := bustime.NewClient("TEST", bustime.CustomBaseURLOption(ts.URL))

	expected := stopsExpectedOutput
	// List of routes to fetch stops for
	routeIDList := []string{"MTA NYCT_M1", "MTA NYCT_M2"}
	actual := client.GetStops(routeIDList...)

	if !reflect.DeepEqual(expected, actual) {
		t.Errorf("bustime.GetStops did not return expected map of routeIDs to stopIDs (expected: %v, received: %v)", expected, actual)
	}
}

var stopsResponses = map[string]string{
	"MTA NYCT_M1": `
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
          "lat": "40.731",
          "lon": "-73.990"
        },
        {
          "id": "MTA_100002",
          "lat": "40.831",
          "lon": "-73.390"
        },
        {
          "id": "MTA_100003",
          "lat": "40.131",
          "lon": "-72.890"
        },
        {
          "id": "MTA_100004",
          "lat": "49.731",
          "lon": "-73.990"
        },
        {
          "id": "MTA_100005",
          "lat": "48.831",
          "lon": "-72.390"
        },
        {
          "id": "MTA_100006",
          "lat": "49.131",
          "lon": "-73.890"
        }
      ]
    }
  }
}
`,
	"MTA NYCT_M2": `
{
  "data": {
    "entry": {
      "stopGroupings": [
        {
          "stopGroups": [
            {
              "id": 0,
              "stopIds": [
                "MTA_200001",
                "MTA_200002",
                "MTA_200003"
              ]
            },
            {
              "id": 1,
              "stopIds": [
                "MTA_200004",
                "MTA_200005",
                "MTA_200006"
              ]
            }
          ]
        }
      ]
    },
    "references": {
      "stops": [
        {
          "id": "MTA_200001",
          "lat": "48.931",
          "lon": "-71.290"
        },
        {
          "id": "MTA_200002",
          "lat": "43.431",
          "lon": "-75.790"
        },
        {
          "id": "MTA_200003",
          "lat": "44.231",
          "lon": "-71.990"
        },
        {
          "id": "MTA_200004",
          "lat": "48.931",
          "lon": "-71.290"
        },
        {
          "id": "MTA_200005",
          "lat": "43.431",
          "lon": "-75.790"
        },
        {
          "id": "MTA_200006",
          "lat": "44.231",
          "lon": "-71.990"
        }
      ]
    }
  }
}
`,
}

var stopsExpectedOutput = map[string]map[int][]bustime.BusStop{
	"MTA NYCT_M1": {
		0: {
			bustime.BusStop{StopID: "MTA_100001", Latitude: 40.731, Longitude: -73.990},
			bustime.BusStop{StopID: "MTA_100002", Latitude: 40.831, Longitude: -73.390},
			bustime.BusStop{StopID: "MTA_100003", Latitude: 40.131, Longitude: -72.890},
		},
		1: {
			bustime.BusStop{StopID: "MTA_100004", Latitude: 49.731, Longitude: -73.990},
			bustime.BusStop{StopID: "MTA_100005", Latitude: 48.831, Longitude: -72.390},
			bustime.BusStop{StopID: "MTA_100006", Latitude: 49.131, Longitude: -73.890},
		},
	},
	"MTA NYCT_M2": {
		0: {
			{StopID: "MTA_200001", Latitude: 48.931, Longitude: -71.290},
			{StopID: "MTA_200002", Latitude: 43.431, Longitude: -75.790},
			{StopID: "MTA_200003", Latitude: 44.231, Longitude: -71.990},
		},
		1: {
			{StopID: "MTA_200004", Latitude: 48.931, Longitude: -71.290},
			{StopID: "MTA_200005", Latitude: 43.431, Longitude: -75.790},
			{StopID: "MTA_200006", Latitude: 44.231, Longitude: -71.990},
		},
	},
}
