package dates

import (
	"log"
	"time"
	"transport/lib/database"
	"transport/lib/stringhelper"
)

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

// SetHour returns a Time with the same year, month and date,
// the specified `newHourValue`, and 0 for minutes, seconds and nanoseconds.
func SetHour(t time.Time, newHourValue int, loc *time.Location) time.Time {
	return time.Date(t.Year(), t.Month(), t.Day(), newHourValue, 0, 0, 0, loc)
}

// Printf can be used to easily log time.Time structs as dates, without the time included.
// Example usage:
//     dates.Printf("Fetching timestamps between %s and %s", time.Now(), time.Now())
// Output:
//     "Fetching timestamps between 2019-05-09 and 2019-05-09"
func Printf(format string, times ...time.Time) {
	dates := make([]string, len(times))
	for i, t := range times {
		dates[i] = t.Format(database.DateFormat)
	}
	log.Printf(format, stringhelper.SliceToInterface(&dates)...)
}
