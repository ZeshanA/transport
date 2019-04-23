package testhelper

import (
	"database/sql"
	"database/sql/driver"
	"log"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
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
func ServeMultiResponseMock(responses map[string]string, extractKeyFromURL func(*url.URL) string) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		keyInMap := extractKeyFromURL(r.URL)
		response := responses[keyInMap]
		_, err := w.Write([]byte(response))
		if err != nil {
			log.Fatal(err)
		}
	}))
}

func ExtractJSONFilepath(fullURL *url.URL) string {
	path := fullURL.Path
	components := strings.Split(path, "/")
	filepath := components[len(components)-1]
	return strings.Replace(filepath, ".json", "", 1)
}

// SetupDBMock can be used to get a sql.DB instance that will
// return the specified rows when queries.
func SetupDBMock(t *testing.T, columnNames []string, rowsToReturn [][]driver.Value, expectedQuery string) (*sql.DB, sqlmock.Sqlmock) {
	// Set up DB mock
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	// Initialise DB mock for rows
	rows := sqlmock.NewRows(columnNames)

	// Assign rows for mock DB to return
	for _, r := range rowsToReturn {
		rows.AddRow(r...)
	}

	// Expect the correct query to be executed
	mock.ExpectQuery(expectedQuery).WillReturnRows(rows)
	return db, mock
}
