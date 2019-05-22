package main

import (
	"fmt"
	"transport/lib/bustime"
	"transport/lib/iohelper"
)

func main() {
	routeID := "MTA NYCT_S78"
	bt := bustime.NewClient(iohelper.GetEnv("MTA_API_KEY"))
	stops := bt.GetStops(routeID)
	fmt.Println(stops)
}
