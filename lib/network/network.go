package network

import (
	"io"
	"log"
	"net/http"
	"os"
	"transport/lib/iohelper"
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
