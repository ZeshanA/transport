package bustime

import (
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"
	"transport/lib/iohelper"

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
	// Create non-nil error so the 'for' loop runs at least once
	err := errors.New("")

	for err != nil {
		resp, err = http.Get(requestURL)
		if err != nil {
			log.Println(err)
			time.Sleep(2 * time.Second)
		}
	}
	defer iohelper.CloseSafely(resp.Body, requestURL)

	rawData, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("fetch.AllAgencies: error parsing list of agencies: %s", err)
	}

	stringData := string(rawData)

	var agencyIDs []string
	gjson.Get(stringData, "data").ForEach(func(_, agency gjson.Result) bool {
		agencyIDs = append(agencyIDs, agency.Get("agency.id").String())
		return true
	})

	return agencyIDs
}
