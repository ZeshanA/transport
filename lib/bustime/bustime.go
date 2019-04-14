package bustime

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"transport/lib/iohelper"
	"transport/lib/network"

	"github.com/avast/retry-go"

	"github.com/tidwall/gjson"
)

const (
	baseURL          = "http://bustime.mta.info/api/where"
	agenciesEndpoint = "agencies-with-coverage.json"
)

type Client struct {
	Key string
}

func (client *Client) GetAgencies() []string {
	URLWithKey := fmt.Sprintf("%s/%s?key=%s", baseURL, agenciesEndpoint, client.Key)
	return client.getAgenciesWithURL(URLWithKey)
}

func (client *Client) getAgenciesWithURL(requestURL string) []string {
	var resp *http.Response

	// Send GET request for agencies; retry a limited number of times if it fails.
	err := retry.Do(network.GetRequestFunc(requestURL, &resp))
	if err != nil {
		log.Fatalf("fetch.AllAgencies: error fetching list of agencies: %s", err)
	}
	defer iohelper.CloseSafely(resp.Body, requestURL)

	// Read response body
	rawData, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("fetch.AllAgencies: error parsing list of agencies: %s", err)
	}

	return client.parseIDsFromAgencyResponse(rawData)
}

func (client *Client) parseIDsFromAgencyResponse(rawResponseBody []byte) []string {
	stringData := string(rawResponseBody)

	var agencyIDs []string
	gjson.Get(stringData, "data").ForEach(func(_, agency gjson.Result) bool {
		agencyIDs = append(agencyIDs, agency.Get("agency.id").String())
		return true
	})

	return agencyIDs
}
