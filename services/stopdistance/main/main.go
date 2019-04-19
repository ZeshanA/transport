package main

import (
	"fmt"
	"log"
	"transport/lib/bustime"
	"transport/lib/database"
	"transport/lib/iohelper"

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
	storeDistances(distances)
}

func storeDistances(distances []stopDistance) {
	database.Store(database.StopDistanceTable, extractStopDistanceColumns, stopDistanceToInterface(distances))
}

func stopDistanceToInterface(distances []stopDistance) []interface{} {
	r := make([]interface{}, len(distances))
	for i, distance := range distances {
		r[i] = distance
	}
	return r
}

func extractStopDistanceColumns(sdEntry interface{}) []interface{} {
	sd := sdEntry.(stopDistance)
	return []interface{}{sd.routeID, sd.fromID, sd.toID, sd.distance, sd.directionID}
}
