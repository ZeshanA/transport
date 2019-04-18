package testhelper

import (
	"log"
	"net/http"
	"net/http/httptest"
	"strings"
)

// ServeMock creates a new httptest server that responds to all requests with
// `response`. Errors are fatal. This function should only be used in test cases.
func ServeMock(response string) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		_, err := w.Write([]byte(response))
		if err != nil {
			log.Fatal(err)
		}
	}))
}

// ServeMultiResponseMock returns a pointer to an httptest.Server that serves a
// response from the `responses` map provided, based on the endpoint of the URL
// in the request, e.g. a request to "/example/endpoint/abc.json" will return
// the string stored under "abc" in the responses map.

// This is intended for use in unit testing. Create a map with the expected
// .JSON endpoints your test will hit, and the corresponding responses that
// should be served, and pass it to ServeMultiResponseMock to spin up
// a test HTTP server that handles all of that for you.
func ServeMultiResponseMock(responses map[string]string, extractKeyFromURL func(string) string) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		keyInMap := extractKeyFromURL(r.URL.Path)
		response := responses[keyInMap]
		_, err := w.Write([]byte(response))
		if err != nil {
			log.Fatal(err)
		}
	}))
}

func ExtractJSONFilepath(url string) string {
	components := strings.Split(url, "/")
	filepath := components[len(components)-1]
	return strings.Replace(filepath, ".json", "", 1)
}
