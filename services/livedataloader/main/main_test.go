package main

import (
	"log"
	"testing"
)

func TestHelloWorld(t *testing.T) {
	expected := "Hello, world!"
	result := getHelloWorld()

	if result != expected {
		log.Fatalf("TestHelloWorld failed: expected %s received %s\n", expected, result)
	}
}
