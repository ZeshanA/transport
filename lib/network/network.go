package network

import (
	"fmt"
	"io"
	"net/http"
	"os"
)

// closeSafely calls an io.ReadCloser's .Close() method
// with error checking, to allow for safe, clean use with `defer`
func closeSafely(item io.ReadCloser, resourcePath string) {
	if err := item.Close(); err != nil {
		fmt.Printf("Error when closing the following resource: %s\n", resourcePath)
	}
}

// DownloadFile will download a URL to a local file. Avoids loading
// the whole file into memory by writing as it receives data
func DownloadFile(URL string, filepath string) error {

	fmt.Printf("Fetching URL %s and saving at %s\n", URL, filepath)

	// Fetch the data
	resp, err := http.Get(URL)
	if err != nil {
		return err
	}
	defer closeSafely(resp.Body, URL)

	// Create the file
	out, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer closeSafely(out, filepath)

	// Write the body to file
	_, err = io.Copy(out, resp.Body)

	return err
}
