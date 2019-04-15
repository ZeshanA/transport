package bustime_test

import (
	"fmt"
	"reflect"
	"testing"
	"transport/lib/bustime"
	"transport/lib/testhelper"
)

func TestClient_GetAgencies(t *testing.T) {
	// Create HTTP server to serve mock JSON response
	agencyList := `[{"agency":{"id":"MTA NYCT"}}, {"agency":{"id":"MTABC"}}]`
	response := fmt.Sprintf(`{"data": %s}`, agencyList)
	ts := testhelper.ServeMock(response)
	defer ts.Close()

	// Create bustime.client
	client := bustime.NewClient("TEST", bustime.CustomBaseURLOption(ts.URL))

	// We expect the following slice of agency ID strings
	expected := []string{"MTA NYCT", "MTABC"}
	actual := *client.GetAgencies()
	if !reflect.DeepEqual(expected, actual) {
		t.Errorf("bustime.GetAgencies did not return expected list of agency IDs (expected: %s, received: %s)", expected, actual)
	}
}

func TestClient_GetAgenciesEmptyResponse(t *testing.T) {
	// Create HTTP server to serve mock JSON response
	response := fmt.Sprintf(`{"data": %s}`, "[]")
	ts := testhelper.ServeMock(response)
	defer ts.Close()

	// Create bustime.client
	client := bustime.NewClient("TEST", bustime.CustomBaseURLOption(ts.URL))

	// We expect nil (no slice should have been allocated, as there were no agency IDs returned)
	var expected []string
	actual := *client.GetAgencies()
	if !reflect.DeepEqual(expected, actual) {
		t.Errorf("bustime.GetAgencies did not return empty list of agency IDs (expected: %s, received: %s)", expected, actual)
	}
}
