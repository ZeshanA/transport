package network

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"transport/lib/iohelper"

	"github.com/avast/retry-go"
)

// DownloadFile will download a URL to a local file. Avoids loading
// the whole file into memory by writing as it receives data.
// From: https://golangcode.com/download-a-file-from-a-url/
func DownloadFile(URL string, filepath string) error {

	log.Printf("Fetching URL %s and saving at %s\n", URL, filepath)

	// Fetch the data
	resp, err := http.Get(URL)
	if err != nil {
		return err
	}
	defer iohelper.CloseSafely(resp.Body, URL)

	// Create the file
	out, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer iohelper.CloseSafely(out, filepath)

	// Write the body to file
	_, err = io.Copy(out, resp.Body)

	return err
}

// GetRequestFunc returns a function that executes an HTTP GET request to
// `requestURL` and changes the pointer stored at responseLocation to point
// to the GET request's response.
// This is primarily intended to be composed with the retry.Do function as follows:
//     var resp *http.Response
//     err := retry.Do(GetRequestFunc(requestURL, &resp))
// This will retry the request 10 times and return an error if none of the 10 requests succeed
func GetRequestFunc(requestURL string, responseLocation **http.Response) func() error {
	return func() error {
		response, err := http.Get(requestURL)
		*responseLocation = response
		return err
	}
}

// GetRequestBody fetches the resource at the given URL and
// returns a pointer to a []byte containing just the body of the response.
// If the request fails, it will be retried up to 10 times before logging
// a fatal error. Failure to read the response bytes results in a fatal error.
func GetRequestBody(requestURL string) string {
	var resp *http.Response

	// Send GET request for agencies; retry a limited number of times if it fails.
	err := retry.Do(GetRequestFunc(requestURL, &resp))
	if err != nil {
		log.Fatalf("fetch.AllAgencies: error fetching list of agencies: %s", err)
	}
	defer iohelper.CloseSafely(resp.Body, requestURL)

	// Read response body
	rawData, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("fetch.AllAgencies: error parsing list of agencies: %s", err)
	}

	return string(rawData)
}

func WriteError(msg string, err error, w http.ResponseWriter) {
	formatted := fmt.Sprintf(msg, err)
	log.Println(formatted)
	w.Write([]byte(formatted))
}
