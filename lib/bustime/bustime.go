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
)

type client struct {
	Key     string
	BaseURL string
}

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

func CustomBaseURLOption(customBaseURL string) func(*client) error {
	return func(client *client) error {
		client.BaseURL = customBaseURL
		return nil
	}
}

/* Agencies */
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
