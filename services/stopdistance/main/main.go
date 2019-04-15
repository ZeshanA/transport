package main

import (
	"transport/lib/bustime"
	"transport/lib/iohelper"
)

func main() {
	client := bustime.Client{Key: iohelper.GetEnv("MTA_API_KEY")}
	fmt.Println(client.GetAgencies())
}
