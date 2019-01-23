package main

import "testing"

func TestHelloWorld(t *testing.T) {
	message := helloWorld()
	expected := "Hello, world!"
	if message != expected {
		t.Errorf("Message was incorrect, got: %s, expected: %s", message, expected)
	}
}
