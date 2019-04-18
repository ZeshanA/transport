package mapping

import (
	"context"
	"fmt"
	"log"

	"googlemaps.github.io/maps"
)

func RoadDistance(mc *maps.Client, fromLat, fromLon, toLat, toLon float64) (distanceInMetres float64) {
	r := &maps.DistanceMatrixRequest{
		Origins:      []string{fmt.Sprintf("%f,%f", fromLat, fromLon)},
		Destinations: []string{fmt.Sprintf("%f,%f", toLat, toLon)},
	}
	distance, err := mc.DistanceMatrix(context.Background(), r)
	if err != nil {
		log.Fatalf("mapping.RoadDistance: error whilst fetching distance matrix response: %s", err)
	}
	return float64(distance.Rows[0].Elements[0].Distance.Meters)
}
