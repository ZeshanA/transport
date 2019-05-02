package main

import (
	"errors"
	"fmt"
	"log"
	"sync"
	"transport/lib/bus"
	"transport/lib/bustime"
	"transport/lib/mapping"
	"transport/services/labeller/stopdistance"

	"googlemaps.github.io/maps"
)

type distanceResponse struct {
	res []bus.StopDistance
	err error
}

func GetDistances(mc *maps.Client, stopDetails map[string]map[int][]bustime.BusStop, existingSDs map[stopdistance.Key]float64) []bus.StopDistance {
	var distances []bus.StopDistance
	mux := &sync.Mutex{}
	fetched, count := make(chan distanceResponse), 0
	for routeID, directionIDs := range stopDetails {
		for directionID, stopsForDirectionID := range directionIDs {
			go getDistancesAlongRoute(mc, routeID, directionID, stopsForDirectionID, fetched, existingSDs)
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

func getDistancesAlongRoute(mc *maps.Client, routeID string, directionID int, stops []bustime.BusStop, fetched chan distanceResponse, existingSDs map[stopdistance.Key]float64) {
	if len(stops) < 2 {
		fetched <- distanceResponse{nil, errors.New("getDistancesAlongRoute: fewer than 2 stops in list provided")}
		return
	}
	log.Printf("Fetching distances for routeID: %s, directionID: %d\n", routeID, directionID)
	dists := make([]bus.StopDistance, len(stops)-1)
	for i := 0; i < len(stops); i++ {
		for j := i + 1; j < len(stops); j++ {
			from, to := stops[i], stops[j]
			k := stopdistance.Key{RouteID: routeID, DirectionID: directionID, FromID: from.ID, ToID: to.ID}
			if _, exists := existingSDs[k]; exists {
				continue
			}
			dists[i] = bus.StopDistance{
				RouteID: routeID, DirectionID: directionID, FromID: from.ID, ToID: to.ID,
				Distance: mapping.RoadDistance(mc, from.Latitude, from.Longitude, to.Latitude, to.Longitude),
			}
			existingSDs[k] = 123.0
		}
	}
	log.Printf("Succesfully fetched distances for routeID: %s, directionID: %d\n", routeID, directionID)
	fetched <- distanceResponse{dists, nil}
}
