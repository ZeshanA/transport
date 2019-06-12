package response

import (
	"detector/request"
	"log"
	"transport/lib/database"
	"transport/lib/iohelper"
)
import "github.com/pusher/pusher-http-go"

var client = pusher.Client{
	AppID:   "790340",
	Key:     "66f6e62226c2a035a177",
	Secret:  iohelper.GetEnv("PUSHER_SECRET"),
	Cluster: "eu",
	Secure:  true,
}

const departureNotification = "departureNotification"

type Notification struct {
	VehicleID            string
	OptimalDepartureTime database.Timestamp
	PredictedArrivalTime database.Timestamp
}

func SendNotification(params request.JourneyParams, notification Notification) {
	err := client.Trigger(params.Channel, departureNotification, notification)
	if err != nil {
		log.Printf("error sending departure notification: %v", err)
	}
}
