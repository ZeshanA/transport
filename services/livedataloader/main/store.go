package main

import (
	"log"
	"transport/lib/database"
)

func store(dataIncoming chan bool) {
	database.OpenDBConnection()
	for {
		<-dataIncoming
		log.Printf("Data incoming!")
	}
}
