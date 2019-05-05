package dates

import "time"

const HoursInDay = 24

func Equal(a time.Time, b time.Time) bool {
	dayA, monthA, yearA := a.Date()
	dayB, monthB, yearB := b.Date()
	return dayA == dayB && monthA == monthB && yearA == yearB
}

// DaysBetween returns the number of days between dates
// `a` and `b`. `b` should be the date that occurs *after* `a`.
// This function is *inclusive* of `a` and `b`:
// e.g. DaysBetween(1st April, 4th April) = 4 days (1st April, 2nd April, 3rd April, 4th April)
func DaysBetween(a time.Time, b time.Time) int {
	return int(b.Sub(a).Hours()/HoursInDay) + 1
}
