package main

import (
	"fmt"
	"transport/lib/bustime"
	"transport/lib/iohelper"
)

func main() {
	client := bustime.NewClient(iohelper.GetEnv("MTA_API_KEY"))
	fmt.Println(client.GetAgencies())
}
