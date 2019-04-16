package bustime

import (
	"fmt"
	"log"
)

const (
	defaultBaseURL    = "http://bustime.mta.info/api/where"
	defaultAPIVersion = "2"
)

type client struct {
	key     string
	baseURL string
	// Query string containing params that *must* be sent with each request,
	// namely the key and API version, e.g. "key=abc&version=2"
	MandatoryParams string
}

// NewClient creates a new bustime.client
// The API `key` parameter is mandatory, the remainder of the
// parameters are optional and are the functions suffixed with
// 'Option' in this file.
// (see: Functional Options pattern – https://dave.cheney.net/2014/10/17/functional-options-for-friendly-apis)

// Example Usage:
// client := bustime.NewClient("API_KEY", CustomBaseURLOption("http://google.com/"))
func NewClient(key string, options ...func(*client) error) *client {
	client := client{key: key, baseURL: defaultBaseURL}
	for _, option := range options {
		err := option(&client)
		if err != nil {
			log.Fatalf("bustime.client initialisation error: %s", err)
		}
	}
	client.MandatoryParams = fmt.Sprintf("key=%s&version=%s", client.key, defaultAPIVersion)
	return &client
}

// CustomBaseURLOption returns a *function* that can be passed
// to the NewClient constructor to override the default base URL
func CustomBaseURLOption(customBaseURL string) func(*client) error {
	return func(client *client) error {
		client.baseURL = customBaseURL
		return nil
	}
}
