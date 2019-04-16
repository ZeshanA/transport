package bustime

import (
	"fmt"
	"log"
	"transport/lib/jsonhelper"
	"transport/lib/network"

	"github.com/tidwall/gjson"
)

const (
	agenciesEndpoint = "agencies-with-coverage.json"
	routesEndpoint   = "routes-for-agency"
	stopsEndpoint    = "stops-for-route"
)

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
// Agencies
////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
func (client *client) GetAgencies() []string {
	URLWithKey := fmt.Sprintf("%s/%s?%s", client.baseURL, agenciesEndpoint, client.MandatoryParams)
	jsonResponse := network.GetRequestBody(URLWithKey)
	return jsonhelper.ExtractNested(jsonResponse, "data.list.#.agencyId")
}

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
// Routes
////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
func (client *client) GetRoutes(agencyIDs ...string) []string {
	var routeIDs []string
	for _, agencyID := range agencyIDs {
		URLWithKey := fmt.Sprintf("%s/%s/%s.json?%s", client.baseURL, routesEndpoint, agencyID, client.MandatoryParams)
		jsonResponse := network.GetRequestBody(URLWithKey)
		newRouteIDs := jsonhelper.ExtractNested(jsonResponse, "data.list.#.id")
		routeIDs = append(routeIDs, newRouteIDs...)
	}
	return routeIDs
}

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
// Stops
////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

// GetStops takes a collection of routeIDs and returns a map of
// the form: routeID -> []stopID
func (client *client) GetStops(routeIDs ...string) map[string][]string {
	routeIDtoStopIDs := map[string][]string{}
	done := make(chan string)
	for _, routeID := range routeIDs {
		go client.getStopsForSingleRoute(routeIDtoStopIDs, routeID, done)
	}
	// Wait for all go routines to report completion to the 'done' channel
	for i := 0; i < len(routeIDs); i++ {
		completedRouteID := <-done
		log.Printf("Succesfully stored stops for route ID: %s\n", completedRouteID)
	}
	return routeIDtoStopIDs
}

func (client *client) getStopsForSingleRoute(routeIDtoStopIDs map[string][]string, routeID string, done chan string) {
	log.Printf("Fetching stops for route ID: %s\n", routeID)
	URLWithKey := fmt.Sprintf(
		"%s/%s/%s.json?%s&includePolylines=false",
		client.baseURL, stopsEndpoint, routeID, client.MandatoryParams,
	)
	// Fetch JSON response containing stopIDs for current routeID
	jsonString := network.GetRequestBody(URLWithKey)
	// Extract stopIDs nested within the response and convert to a []string
	stopIDs := gjson.Get(jsonString, "data.entry.stopIds").Array()
	stringStopIDs := jsonhelper.ResultArrayToStringArray(stopIDs)
	// Store list of stopIDs under given routeID
	routeIDtoStopIDs[routeID] = stringStopIDs
	// Write routeID to channel to mark as completed
	done <- routeID
}
