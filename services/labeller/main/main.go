package main

import (
	"fmt"
	"labeller/stopdistance"
	"transport/lib/database"
)

var db = database.OpenDBConnection()

func main() {
	stopDistances := stopdistance.Get(db)
	fmt.Println(stopDistances)
}
