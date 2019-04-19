package dates

import "time"

func Equal(a time.Time, b time.Time) bool {
	dayA, monthA, yearA := a.Date()
	dayB, monthB, yearB := b.Date()
	return dayA == dayB && monthA == monthB && yearA == yearB
}
