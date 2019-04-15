package testhelper

import (
	"log"
	"net/http"
	"net/http/httptest"
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
