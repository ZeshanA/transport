package main

import "testing"

func TestHelloWorld(t *testing.T) {
	message := helloWorld()
	expected := "Hello, world!"
	if message != expected {
		printMismatchError(t, message, expected)
	}
}

func printMismatchError(t *testing.T, actual string, expected string) {
	t.Errorf("Message was incorrect, got: %s, expected: %s", actual, expected)
}
