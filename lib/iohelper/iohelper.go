package iohelper

import (
	"flag"
	"fmt"
	"io"
	"log"
	"os"
)

// CloseSafely calls an io.ReadCloser's .Close() method
// with error checking, to allow for safe, clean use with `defer`
func CloseSafely(item io.ReadCloser, resourcePath string) {
	if err := item.Close(); err != nil {
		fmt.Printf("Error when closing the following resource: %s\n", resourcePath)
	}
}

// GetEnv fetches the value of the environment variable stored
// at key. If the environment variable isn't set, the program
// will exit(1).
func GetEnv(key string) string {
	value := os.Getenv(key)
	if value == "" && flag.Lookup("test.v") == nil {
		log.Fatalf("%s not set\n", value)
	}
	return value
}
