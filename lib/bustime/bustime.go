package bustime

import (
	"fmt"
	"log"
	"transport/lib/network"

	"github.com/tidwall/gjson"
)

/* TODO: Optimise all client methods to return pointers to avoid
   copying return values around. */

const (
	defaultBaseURL   = "http://bustime.mta.info/api/where"
	agenciesEndpoint = "agencies-with-coverage.json"
	routesEndpoint   = "routes-for-agency"
)

type client struct {
	Key     string
	BaseURL string
}

// NewClient creates a new bustime.client
// The API `key` parameter is mandatory, the remainder of the
// parameters are optional and are the functions suffixed with
// 'Option' in this file.

// Example Usage:
// client := bustime.NewClient("API_KEY", CustomBaseURLOption("http://google.com/"))
func NewClient(key string, options ...func(*client) error) *client {
	client := client{Key: key, BaseURL: defaultBaseURL}
	for _, option := range options {
		err := option(&client)
		if err != nil {
			log.Fatalf("bustime.client initialisation error: %s", err)
		}
	}
	return &client
}

// CustomBaseURLOption returns a *function* that can be passed
// to the NewClient constructor to override the default base URL
func CustomBaseURLOption(customBaseURL string) func(*client) error {
	return func(client *client) error {
		client.BaseURL = customBaseURL
		return nil
	}
}

// Agencies
func (client *client) GetAgencies() *[]string {
	URLWithKey := fmt.Sprintf("%s/%s?key=%s", client.BaseURL, agenciesEndpoint, client.Key)
	rawData := network.GetRequestBody(URLWithKey)
	return client.parseIDsFromAgencyResponse(rawData)
}

func (client *client) parseIDsFromAgencyResponse(rawResponseBody *[]byte) *[]string {
	stringData := string(*rawResponseBody)

	var agencyIDs []string
	gjson.Get(stringData, "data").ForEach(func(_, agency gjson.Result) bool {
		agencyIDs = append(agencyIDs, agency.Get("agency.id").String())
		return true
	})

	return &agencyIDs
}

// Routes
func (client *client) GetRoutes(agencyIDs ...string) []string {
	var routeIDs []string
	for _, agencyID := range agencyIDs {
		URLWithKey := fmt.Sprintf("%s/%s/%s.json?key=%s", client.BaseURL, routesEndpoint, agencyID, client.Key)
		rawData := network.GetRequestBody(URLWithKey)
		routeIDs = append(routeIDs, client.parseIDsFromRoutesResponse(rawData)...)
	}
	return routeIDs
}

func (client *client) parseIDsFromRoutesResponse(rawResponseBody *[]byte) []string {
	stringData := string(*rawResponseBody)
	routeObjects := gjson.Get(stringData, "data.list").Array()
	routeIDs := make([]string, len(routeObjects))
	for i, route := range routeObjects {
		routeIDs[i] = route.Get("id").String()
	}
	return routeIDs
}
