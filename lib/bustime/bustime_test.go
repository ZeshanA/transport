package bustime_test

import (
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
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

func TestClient_GetRoutesWithURL(t *testing.T) {
	// The agency IDs that we want the routes for
	agencyIDsList := []string{"MTA NYCT", "MTABC"}

	// Construct mocked JSON responses containing routes for each agency ID
	responses := map[string]string{
		"MTA NYCT": `{"data": {"list": [{"id": "MTA NYCT_BX4A"}, {"id": "MTA NYCT_M31"}]}}`,
		"MTABC":    `{"data": {"list": [{"id": "MTABC_BX4A"}, {"id": "MTABC_M31"}]}}`,
	}

	// Create HTTP server to serve the correct JSON response, based on the agency ID in the request URL
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var response string
		if strings.Contains(r.URL.Path, "MTA NYCT") {
			response = responses["MTA NYCT"]
		} else {
			response = responses["MTABC"]
		}
		_, err := w.Write([]byte(response))
		if err != nil {
			log.Fatal(err)
		}
	}))
	defer ts.Close()

	// Create bustime.client
	client := bustime.NewClient("TEST", bustime.CustomBaseURLOption(ts.URL))

	expected := []string{"MTA NYCT_BX4A", "MTA NYCT_M31", "MTABC_BX4A", "MTABC_M31"}
	actual := client.GetRoutes(agencyIDsList...)
	if !reflect.DeepEqual(expected, actual) {
		t.Errorf("bustime.GetRoutes did not return expected list of route IDs (expected: %s, received: %s)", expected, actual)
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
		t.Errorf("bustime.GetRoutes did not return empty list of route IDs (expected: %s, received: %s)", expected, actual)
	}
}
