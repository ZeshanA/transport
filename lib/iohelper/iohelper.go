package iohelper

import (
	"fmt"
	"io"
)

// CloseSafely calls an io.ReadCloser's .Close() method
// with error checking, to allow for safe, clean use with `defer`
func CloseSafely(item io.ReadCloser, resourcePath string) {
	if err := item.Close(); err != nil {
		fmt.Printf("Error when closing the following resource: %s\n", resourcePath)
	}
}
