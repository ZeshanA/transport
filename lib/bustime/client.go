package bustime

import (
	"fmt"
	"log"
)

const (
	defaultBaseURL    = "http://bustime.mta.info/api/where"
	defaultAPIVersion = "2"
)

type Client struct {
	key     string
	baseURL string
	// Query string containing params that *must* be sent with each request,
	// namely the key and API version, e.g. "key=abc&version=2"
	MandatoryParams string
}

// NewClient creates a new bustime.Client
// The API `key` parameter is mandatory, the remainder of the
// parameters are optional and are the functions suffixed with
// 'Option' in this file.
// (see: Functional Options pattern – https://dave.cheney.net/2014/10/17/functional-options-for-friendly-apis)

// Example Usage:
// Client := bustime.NewClient("API_KEY", CustomBaseURLOption("http://google.com/"))
func NewClient(key string, options ...func(*Client) error) *Client {
	client := Client{key: key, baseURL: defaultBaseURL}
	for _, option := range options {
		err := option(&client)
		if err != nil {
			log.Fatalf("bustime.Client initialisation error: %s", err)
		}
	}
	client.MandatoryParams = fmt.Sprintf("key=%s&version=%s", client.key, defaultAPIVersion)
	return &client
}

// CustomBaseURLOption returns a *function* that can be passed
// to the NewClient constructor to override the default base URL
func CustomBaseURLOption(customBaseURL string) func(*Client) error {
	return func(client *Client) error {
		client.baseURL = customBaseURL
		return nil
	}
}
