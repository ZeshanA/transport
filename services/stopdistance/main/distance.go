package main

import (
	"errors"
	"fmt"
	"log"
	"transport/lib/bustime"
	"transport/lib/mapping"

	"googlemaps.github.io/maps"
)

func GetDistances(mc *maps.Client, stopDetails map[string]map[int][]bustime.BusStop) []stopDistance {
	var distances []stopDistance
	for routeID, directionIDs := range stopDetails {
		for directionID, stopsForDirectionID := range directionIDs {
			distsForRoute, err := getDistancesAlongRoute(mc, routeID, directionID, stopsForDirectionID)
			if err != nil {
				fmt.Println(err)
			}
			distances = append(distances, distsForRoute...)
		}
	}
	return distances
}

func getDistancesAlongRoute(mc *maps.Client, routeID string, directionID int, stops []bustime.BusStop) ([]stopDistance, error) {
	if len(stops) < 2 {
		return nil, errors.New("getDistancesAlongRoute: fewer than 2 stops in list provided")
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
	log.Printf("Succesfully fetched distances for routeID: %s\n, directionID: %d", routeID, directionID)
	return dists, nil
}
