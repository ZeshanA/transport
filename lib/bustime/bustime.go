package bustime

import (
	"fmt"
	"transport/lib/network"

	"github.com/tidwall/gjson"
)

const (
	baseURL          = "http://bustime.mta.info/api/where"
	agenciesEndpoint = "agencies-with-coverage.json"
)

type Client struct {
	Key string
}

/* Agencies */
func (client *Client) GetAgencies() []string {
	URLWithKey := fmt.Sprintf("%s/%s?key=%s", baseURL, agenciesEndpoint, client.Key)
	return client.getAgenciesWithURL(URLWithKey)
}

func (client *Client) getAgenciesWithURL(requestURL string) []string {
	rawData := network.GetRequestBody(requestURL)
	return client.parseIDsFromAgencyResponse(rawData)
}

func (client *Client) parseIDsFromAgencyResponse(rawResponseBody *[]byte) []string {
	stringData := string(*rawResponseBody)

	var agencyIDs []string
	gjson.Get(stringData, "data").ForEach(func(_, agency gjson.Result) bool {
		agencyIDs = append(agencyIDs, agency.Get("agency.id").String())
		return true
	})

	return agencyIDs
}
