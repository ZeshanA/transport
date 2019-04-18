package mapping

import (
	"fmt"
	"log"
	"transport/lib/network"
	"transport/lib/urlhelper"

	"github.com/tidwall/gjson"
)

const (
	defaultBaseURL           = "https://api.mapbox.com"
	directionsMatrixEndpoint = "directions-matrix/v1"
)

type client struct {
	key     string
	baseURL string
}

// NewClient creates a new mapping.client
// The API `key` parameter is mandatory, the remainder of the
// parameters are optional and are the functions suffixed with
// 'Option' in this file.
// (see: Functional Options pattern – https://dave.cheney.net/2014/10/17/functional-options-for-friendly-apis)

// Example Usage:
// client := mapping.NewClient("API_KEY", CustomBaseURLOption("http://google.com/"))
func NewClient(key string, options ...func(*client) error) *client {
	client := client{key: key, baseURL: defaultBaseURL}
	for _, option := range options {
		err := option(&client)
		if err != nil {
			log.Fatalf("mapping.client initialisation error: %s", err)
		}
	}
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

func (c *client) RoadDistance(fromLat, fromLon, toLat, toLon float64) (distanceInMetres float64) {
	URL := c.constructRoadDistanceURL(fromLat, fromLon, toLat, toLon)
	resp := network.GetRequestBody(URL)
	parsedResp := gjson.Parse(resp)
	if parsedResp.Get("code").String() != "Ok" {
		log.Panicf("mapping: error reported in MapBox API response: %s", resp)
	}
	return parsedResp.Get("distances.0.0").Float()
}

func (c *client) constructRoadDistanceURL(fromLat, fromLon, toLat, toLon float64) string {
	coords := fmt.Sprintf("%f,%f;%f,%f", fromLat, fromLon, toLat, toLon)
	rawURL := fmt.Sprintf("%s/%s/mapbox/driving/%s", c.baseURL, directionsMatrixEndpoint, coords)
	params := map[string]string{
		"access_token": c.key,
		"sources":      "0",
		"destinations": "1",
		"annotations":  "distance",
	}
	return fmt.Sprintf("%s%s", rawURL, urlhelper.BuildQueryString(params))
}
