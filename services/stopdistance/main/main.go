package main

import (
	"log"
	"transport/lib/bus"
	"transport/lib/bustime"
	"transport/lib/database"
	"transport/lib/iohelper"

	"googlemaps.github.io/maps"
)

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
	stopDetails := bt.GetStops(routes...)

	// Calculate distances between stops and store in DB
	distances := GetDistances(mc, stopDetails)
	storeDistances(distances)
}

func storeDistances(distances []bus.StopDistance) {
	database.Store(database.StopDistanceTable, extractStopDistanceColumns, stopDistanceToInterface(distances))
}

func stopDistanceToInterface(distances []bus.StopDistance) []interface{} {
	r := make([]interface{}, len(distances))
	for i, distance := range distances {
		r[i] = distance
	}
	return r
}

func extractStopDistanceColumns(sdEntry interface{}) []interface{} {
	sd := sdEntry.(bus.StopDistance)
	return []interface{}{sd.RouteID, sd.FromID, sd.ToID, sd.Distance, sd.DirectionID}
}
