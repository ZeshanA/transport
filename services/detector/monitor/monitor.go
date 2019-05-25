package monitor

import (
	"detector/request"
	"time"
	"transport/lib/bustime"
	"transport/lib/sleep"
)

// Start monitoring 3 journey times before the desired arrival time
const monitoringWindowSize = 5

func LiveBuses(avgTime int, predictedTime int, params request.JourneyParams, stopList []bustime.BusStop) {
	// Sleep until the monitoring window starts (a few hours before the arrival time)
	sleepUntilMonitoringWindow(predictedTime, params)
}

func sleepUntilMonitoringWindow(predictedTime int, params request.JourneyParams) {
	predictedDuration := time.Duration(predictedTime) * time.Second
	sleep.Until(params.ArrivalTime.Add(-monitoringWindowSize * predictedDuration))
}
