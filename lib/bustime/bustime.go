package bustime

import (
	"fmt"
	"log"
	"transport/lib/network"

	"github.com/tidwall/gjson"
)

const (
	defaultBaseURL    = "http://bustime.mta.info/api/where"
	defaultAPIVersion = "2"
	agenciesEndpoint  = "agencies-with-coverage.json"
	routesEndpoint    = "routes-for-agency"
)

type client struct {
	Key             string
	BaseURL         string
	MandatoryParams string
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
	client.MandatoryParams = fmt.Sprintf("key=%s&version=%s", client.Key, defaultAPIVersion)
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

// parseIDsFromResponse can be used to extract an array of string IDs
// from a JSON response.
// rawResponseBody: pointer to an []byte containing a JSON string
// pathToArray: a string path to the array of objects that you wish to extract IDs
//              from, within the JSON string (e.g. "data.routes")
// nameOfIDField: the name of the ID field within each object in the JSON array (e.g. "vehicleID")
func (client *client) parseIDsFromResponse(rawResponseBody *[]byte, pathToArray string, nameOfIDField string) *[]string {
	stringData := string(*rawResponseBody)

	var listOfIDs []string
	gjson.Get(stringData, pathToArray).ForEach(func(_, agency gjson.Result) bool {
		listOfIDs = append(listOfIDs, agency.Get(nameOfIDField).String())
		return true
	})

	return &listOfIDs
}

// Agencies
func (client *client) GetAgencies() *[]string {
	URLWithKey := fmt.Sprintf("%s/%s?%s", client.BaseURL, agenciesEndpoint, client.MandatoryParams)
	rawData := network.GetRequestBody(URLWithKey)
	return client.parseIDsFromResponse(rawData, "data.list", "agencyId")
}

// Routes
func (client *client) GetRoutes(agencyIDs ...string) *[]string {
	var routeIDs []string
	for _, agencyID := range agencyIDs {
		URLWithKey := fmt.Sprintf("%s/%s/%s.json?%s", client.BaseURL, routesEndpoint, agencyID, client.MandatoryParams)
		rawData := network.GetRequestBody(URLWithKey)
		routeIDs = append(routeIDs, *client.parseIDsFromResponse(rawData, "data.list", "id")...)
	}
	return &routeIDs
}
