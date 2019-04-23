package main

import (
	"fmt"
	"labeller/stopdistance"
	"transport/lib/database"
)

var dbConn = database.OpenDBConnection()

func main() {
	stopDistances := stopdistance.Get(dbConn)
	fmt.Println(stopDistances)
}
