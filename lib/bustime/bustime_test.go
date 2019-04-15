package bustime_test

import (
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
	"transport/lib/bustime"
)

// TODO: Add more than one test case for each method

const agencyList = `[{"agency":{"id":"MTA NYCT"}}, {"agency":{"id":"MTABC"}}]`

func TestClient_GetAgencies(t *testing.T) {
	// Construct mocked JSON response
	response := fmt.Sprintf(`{"code":200, "currentTime":1555241357581,"data": %s}`, agencyList)

	// Create HTTP server to serve mock JSON response
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		_, err := w.Write([]byte(response))
		if err != nil {
			log.Fatal(err)
		}
	}))
	defer ts.Close()

	// Create bustime.client
	client := bustime.NewClient("TEST", bustime.CustomBaseURLOption(ts.URL))

	expected := []string{"MTA NYCT", "MTABC"}
	actual := *client.GetAgencies()
	if !reflect.DeepEqual(expected, actual) {
		t.Errorf("bustime.GetAgencies did not return expected list of agency IDs (expected: %s, received: %s)", expected, actual)
	}
}
