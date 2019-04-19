package main

import (
	"errors"
	"fmt"
	"log"
	"sync"
	"transport/lib/bustime"
	"transport/lib/mapping"

	"googlemaps.github.io/maps"
)

type distanceResponse struct {
	res []stopDistance
	err error
}

func GetDistances(mc *maps.Client, stopDetails map[string]map[int][]bustime.BusStop) []stopDistance {
	var distances []stopDistance
	mux := &sync.Mutex{}
	fetched, count := make(chan distanceResponse), 0
	for routeID, directionIDs := range stopDetails {
		for directionID, stopsForDirectionID := range directionIDs {
			go getDistancesAlongRoute(mc, routeID, directionID, stopsForDirectionID, fetched)
			count++
		}
	}
	for i := 0; i < count; i++ {
		distResp := <-fetched
		if distResp.err != nil {
			fmt.Println(distResp.err)
		} else {
			mux.Lock()
			distances = append(distances, distResp.res...)
			mux.Unlock()
		}
	}
	return distances
}

func getDistancesAlongRoute(mc *maps.Client, routeID string, directionID int, stops []bustime.BusStop, fetched chan distanceResponse) {
	if len(stops) < 2 {
		fetched <- distanceResponse{nil, errors.New("getDistancesAlongRoute: fewer than 2 stops in list provided")}
		return
	}
	log.Printf("Fetching distances for routeID: %s, directionID: %d\n", routeID, directionID)
	dists := make([]stopDistance, len(stops)-1)
	for i, j := 0, 1; j < len(stops); i, j = i+1, j+1 {
		from, to := stops[i], stops[j]
		dists[i] = stopDistance{
			routeID: routeID, directionID: directionID, fromID: from.ID, toID: to.ID,
			distance: mapping.RoadDistance(mc, from.Latitude, from.Longitude, to.Latitude, to.Longitude),
		}
	}
	log.Printf("Succesfully fetched distances for routeID: %s, directionID: %d\n", routeID, directionID)
	fetched <- distanceResponse{dists, nil}
}
