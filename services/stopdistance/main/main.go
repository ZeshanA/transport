package main

import (
	"fmt"
	"transport/lib/bustime"
	"transport/lib/iohelper"
)

func main() {
	client := bustime.NewClient(iohelper.GetEnv("MTA_API_KEY"))
	agencies := client.GetAgencies()
	routes := client.GetRoutes(*agencies...)
	fmt.Println(len(routes))
}
