package eval

import (
	"detector/request"
	"fmt"
	"log"
	"math/rand"
	"time"
	"transport/lib/bustime"
	"transport/lib/database"
	"transport/lib/iohelper"
	"transport/lib/math"
)

func Evaluate() {
	log.Println("Evaluation mode")
	params := generateRandomParams()
	fmt.Println(params)
}

func generateRandomParams() request.JourneyParams {
	// Get list of all routeIDs
	bt := bustime.NewClient(iohelper.GetEnv("MTA_API_KEY"))
	routeIDs := bt.GetRoutes("MTA NYCT", "MTABC")
	// Initialise global pseudo-RNG
	rand.Seed(time.Now().Unix())
	// Select a random routeID and directionID
	randomRouteID := routeIDs[rand.Intn(len(routeIDs))]
	possibleDirections := 2
	randomDirectionID := rand.Intn(possibleDirections)
	// Fetch stops for the given routeID and directionID
	stops := bt.GetStops(randomRouteID)[randomRouteID][randomDirectionID]
	// Select a random source stop (excluding the final stop on the route, as there are
	// no possible journeys in that case)
	sourceStopIndex := rand.Intn(len(stops) - 1)
	sourceStop := stops[sourceStopIndex].ID
	destStop := stops[math.RandInRange(sourceStopIndex+1, len(stops))].ID
	// Pick random arrival time
	delay := time.Duration(math.RandInRange(1, 5)) * time.Hour
	arrivalTime := time.Now().In(database.TimeLoc).Add(delay)
	return request.JourneyParams{
		RouteID: randomRouteID, DirectionID: randomDirectionID,
		FromStop: sourceStop, ToStop: destStop,
		ArrivalTime: arrivalTime,
	}
}
