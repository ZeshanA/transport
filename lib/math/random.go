package math

import (
	"math/rand"
	"time"
)

// RandInRange returns an integer in the range [min, max)
func RandInRange(min, max int) int {
	rand.Seed(time.Now().Unix())
	return rand.Intn(max-min) + min
}
