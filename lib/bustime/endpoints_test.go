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
	responses := map[string]string{
		"MTA NYCT_M1": `{"data":{"entry":{"stopIds":["MTA_100001","MTA_100002","MTA_100003"]}}}`,
		"MTA NYCT_M2": `{"data":{"entry":{"stopIds":["MTA_200001","MTA_200002","MTA_200003"]}}}`,
	}
	ts := testhelper.ServeMultiResponseMock(responses)
	defer ts.Close()

	// Create bustime.client
	client := bustime.NewClient("TEST", bustime.CustomBaseURLOption(ts.URL))

	expected := map[string][]string{
		"MTA NYCT_M1": {"MTA_100001", "MTA_100002", "MTA_100003"},
		"MTA NYCT_M2": {"MTA_200001", "MTA_200002", "MTA_200003"},
	}
	// List of routes to fetch stops for
	routeIDList := []string{"MTA NYCT_M1", "MTA NYCT_M2"}
	actual := client.GetStops(routeIDList...)

	if !reflect.DeepEqual(expected, actual) {
		t.Errorf("bustime.GetStops did not return expected map of routeIDs to stopIDs (expected: %s, received: %s)", expected, actual)
	}
}
