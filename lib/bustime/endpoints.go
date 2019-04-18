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

type BusStop struct {
	ID        string
	Latitude  float64
	Longitude float64
}

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
// the form: routeID -> directionID -> []stopID
func (client *client) GetStops(routeIDs ...string) map[string]map[int][]BusStop {
	mapOfStops := map[string]map[int][]BusStop{}
	done := make(chan string)
	for _, routeID := range routeIDs {
		go client.populateStopsForRoute(mapOfStops, routeID, done)
	}
	// Wait for all go routines to report completion to the 'done' channel
	for i := 0; i < len(routeIDs); i++ {
		completedRouteID := <-done
		log.Printf("Succesfully stored stops for route ID: %s\n", completedRouteID)
	}
	return mapOfStops
}

func (client *client) populateStopsForRoute(mapOfStops map[string]map[int][]BusStop, routeID string, done chan string) {
	log.Printf("Fetching stops for route ID: %s\n", routeID)
	// Initialise inner map for this routeID
	mapOfStops[routeID] = map[int][]BusStop{}
	// Construct the URL to fetch data from
	URLWithKey := fmt.Sprintf(
		"%s/%s/%s.json?%s&includePolylines=false",
		client.baseURL, stopsEndpoint, routeID, client.MandatoryParams,
	)
	// Fetch JSON response containing stopIDs for current routeID
	jsonString := network.GetRequestBody(URLWithKey)
	// Get the list of travel directions for this routeID
	// travelDirections := gjson.Get(jsonString, "data.entry.stopGroupings[0].stopGroups").Array()
	travelDirections := gjson.Get(jsonString, "data.entry.stopGroupings").Array()[0].Get("stopGroups").Array()
	// Get a map of stopIDs -> stopDetails (JSON strings which contain all of the lat/lon details for each stopID)
	stopDetails := getStopDetails(jsonString)
	// For the current routeID, populate each direction (0 or 1) with BusStop structs
	for _, direction := range travelDirections {
		client.populateDirectionWithStops(mapOfStops, stopDetails, routeID, direction)
	}
	// Write to channel to mark the current routeID as completed
	done <- routeID
}

func (client *client) populateDirectionWithStops(
	mapOfStops map[string]map[int][]BusStop,
	stopDetails map[string]gjson.Result,
	routeID string, direction gjson.Result) {
	// Convert direction to integer (0 or 1)
	directionID := int(direction.Get("id").Int())
	// Extract list of stopIDs from JSON
	stopIDs := jsonhelper.ResultArrayToStringArray(direction.Get("stopIds").Array())
	// Initialise the list of BusStop structs under the current routeID and directionID
	mapOfStops[routeID][directionID] = make([]BusStop, len(stopIDs))
	// Construct a BusStop struct for each stop and store it in the map
	for i, id := range stopIDs {
		curStopDetails := stopDetails[id]
		lat, lon := curStopDetails.Get("lat").Float(), curStopDetails.Get("lon").Float()
		stopStruct := BusStop{ID: id, Latitude: lat, Longitude: lon}
		mapOfStops[routeID][directionID][i] = stopStruct
	}
}

// Constructs a map of stopID -> stopDetails from the JSON.
// The stop details are in an array in the JSON, converting to a map
// keyed by stopID allows for O(1) extraction of lat/lon values given a stopID.
func getStopDetails(jsonString string) map[string]gjson.Result {
	stopIDsMap := map[string]gjson.Result{}
	stops := gjson.Get(jsonString, "data.references.stops").Array()
	for _, stop := range stops {
		stopID := stop.Get("id").String()
		stopIDsMap[stopID] = stop
	}
	return stopIDsMap
}
