package sleep

import (
	"log"
	"time"
	"transport/lib/database"
)

func Until(wakeTime time.Time) {
	// Sleep until it's 4am
	t := time.Now().In(database.TimeLoc)
	timeToSleep := wakeTime.Sub(t)
	log.Printf("Sleeping until %s\n", wakeTime.Format(database.TimeFormat))
	time.Sleep(timeToSleep)
}
