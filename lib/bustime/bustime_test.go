package bustime

import (
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

func TestClient_GetAgenciesWithURL(t *testing.T) {
	// Construct mocked JSON response
	const agencyList = `[{"agency":{"id":"MTA NYCT"}}, {"agency":{"id":"MTABC"}}]`
	response := fmt.Sprintf(`{"code":200, "currentTime":1555241357581,"data": %s}`, agencyList)

	// Create HTTP server to serve mock JSON response
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		_, err := w.Write([]byte(response))
		if err != nil {
			log.Fatal(err)
		}
	}))
	defer ts.Close()

	// Create bustime.Client
	client := Client{Key: "TEST"}

	expected := []string{"MTA NYCT", "MTABC"}
	actual := client.getAgenciesWithURL(ts.URL)
	if !reflect.DeepEqual(expected, actual) {
		t.Errorf("bustime.GetAgencies did not return expected list of agency IDs (expected: %s, received: %s)", expected, actual)
	}
}
