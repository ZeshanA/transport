package main

import (
	"errors"
	"fmt"
	"log"
	"transport/lib/bustime"
	"transport/lib/iohelper"
	"transport/lib/mapping"

	"googlemaps.github.io/maps"
)

type stopDistance struct {
	routeID     string
	directionID int
	fromID      string
	toID        string
	distance    float64
}

func (sd *stopDistance) String() string {
	return fmt.Sprintf("%s – Direction %d – From %s – To %s = %f metres", sd.routeID, sd.directionID, sd.fromID, sd.toID, sd.distance)
}

func main() {
	bt := bustime.NewClient(iohelper.GetEnv("MTA_API_KEY"))
	mc, err := maps.NewClient(maps.WithAPIKey(iohelper.GetEnv("GOOGLE_MAPS_API_KEY")))
	if err != nil {
		log.Panicf("main: failed to initialise Maps API client: %s", err)
	}

	// Get stopDetails in map with format routeID -> directionID -> []BusStop
	agencies := bt.GetAgencies()
	log.Printf("%d agencies fetched\n", len(agencies))
	routes := bt.GetRoutes(agencies...)
	log.Printf("%d routes fetched\n", len(agencies))
	stopDetails := bt.GetStops(routes[0])

	// Calculate distances between stops and store in DB
	distances := GetDistances(mc, stopDetails)
	StoreDistances(distances)
}

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
	dists := make([]stopDistance, len(stops)-1)
	for i, j := 0, 1; j < len(stops); i, j = i+1, j+1 {
		from, to := stops[i], stops[j]
		dists[i] = stopDistance{
			routeID: routeID, directionID: directionID, fromID: from.ID, toID: to.ID,
			distance: mapping.RoadDistance(mc, from.Latitude, from.Longitude, to.Latitude, to.Longitude),
		}
	}
	return dists, nil
}

// TODO: Complete
func StoreDistances(distances []stopDistance) {

}
