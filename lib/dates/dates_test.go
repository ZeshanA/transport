package dates

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestDatesEqual(t *testing.T) {
	actual := make([]bool, len(datesA))
	for i, dateA := range datesA {
		dateB := datesB[i]
		actual[i] = Equal(dateA, dateB)
	}
	assert.Equal(t, expectedEqualityOutput, actual)
}

var datesA = []time.Time{
	time.Date(2019, 01, 02, 10, 05, 0, 0, time.UTC),
	time.Date(2019, 02, 03, 00, 01, 0, 0, time.UTC),
	time.Date(2020, 03, 04, 00, 01, 0, 0, time.UTC),
	time.Date(2018, 02, 03, 00, 01, 0, 0, time.UTC),
}

var datesB = []time.Time{
	time.Date(2019, 01, 02, 10, 04, 0, 0, time.UTC),
	time.Date(2018, 02, 03, 00, 01, 0, 0, time.UTC),
	time.Date(2020, 03, 04, 00, 01, 0, 0, time.UTC),
	time.Date(2018, 02, 04, 00, 01, 0, 0, time.UTC),
}

var expectedEqualityOutput = []bool{true, false, true, false}
