package network

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"transport/lib/iohelper"
)

// DownloadFile will download a URL to a local file. Avoids loading
// the whole file into memory by writing as it receives data.
// From: https://golangcode.com/download-a-file-from-a-url/
func DownloadFile(URL string, filepath string) error {

	fmt.Printf("Fetching URL %s and saving at %s\n", URL, filepath)

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
