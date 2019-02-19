package main

import "log"

func store(dataIncoming chan bool) {
	for {
		<-dataIncoming
		log.Printf("Data incoming!")
	}
}
