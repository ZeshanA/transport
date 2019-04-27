package labels_test

import (
	"labeller/labels"
	"labeller/stopdistance"
	"testing"
	"time"
	"transport/lib/bus"
	"transport/lib/database"
	"transport/lib/nulltypes"

	"github.com/stretchr/testify/assert"
	"gopkg.in/guregu/null.v3"
)

// Perfect stop case: bus approaches a stop and the distanceToStop goes exactly to 0.
func TestReachesStop(t *testing.T) {
	reachesStopSequence := map[bus.DirectedRoute][]bus.VehicleJourney{
		{"M55", 0, "ABC"}: {
			{
				DistanceFromStop: null.IntFrom(200), StopPointRef: null.StringFrom("1"),
				Timestamp: nulltypes.TimestampFrom(database.Timestamp{Time: time.Date(2019, 04, 23, 16, 30, 00, 0, time.UTC)}),
			},
			{
				DistanceFromStop: null.IntFrom(100), StopPointRef: null.StringFrom("1"),
				Timestamp: nulltypes.TimestampFrom(database.Timestamp{Time: time.Date(2019, 04, 23, 16, 31, 00, 0, time.UTC)}),
			},
			{
				DistanceFromStop: null.IntFrom(50), StopPointRef: null.StringFrom("1"),
				Timestamp: nulltypes.TimestampFrom(database.Timestamp{Time: time.Date(2019, 04, 23, 16, 32, 00, 0, time.UTC)}),
			},
			{
				DistanceFromStop: null.IntFrom(0), StopPointRef: null.StringFrom("1"),
				Timestamp: nulltypes.TimestampFrom(database.Timestamp{Time: time.Date(2019, 04, 23, 16, 33, 00, 0, time.UTC)}),
			},
			{
				DistanceFromStop: null.IntFrom(0), StopPointRef: null.StringFrom("1"),
				Timestamp: nulltypes.TimestampFrom(database.Timestamp{Time: time.Date(2019, 04, 23, 16, 34, 00, 0, time.UTC)}),
			},
		},
	}
	actual := labels.Create(reachesStopSequence, nil)
	expected := getExpectedLabelledJourneys(reachesStopSequence, []int{180, 120, 60})
	assert.Equal(t, expected, actual)
}

// Stop Not Reached case: the distance to stop never goes to 0 and the stopID never changes.
// We shouldn't get any results in this case, we can't know when it got to the stop.
func TestStopNotReached(t *testing.T) {
	stopNotReached := map[bus.DirectedRoute][]bus.VehicleJourney{
		{"M55", 0, "ABC"}: {
			{
				DistanceFromStop: null.IntFrom(200), StopPointRef: null.StringFrom("1"),
				Timestamp: nulltypes.TimestampFrom(database.Timestamp{Time: time.Date(2019, 04, 23, 16, 35, 00, 0, time.UTC)}),
			},
			{
				DistanceFromStop: null.IntFrom(100), StopPointRef: null.StringFrom("1"),
				Timestamp: nulltypes.TimestampFrom(database.Timestamp{Time: time.Date(2019, 04, 23, 16, 34, 00, 0, time.UTC)}),
			},
			{
				DistanceFromStop: null.IntFrom(50), StopPointRef: null.StringFrom("1"),
				Timestamp: nulltypes.TimestampFrom(database.Timestamp{Time: time.Date(2019, 04, 23, 16, 32, 00, 0, time.UTC)}),
			},
		},
	}
	expected := getExpectedLabelledJourneys(stopNotReached, []int{})
	actual := labels.Create(stopNotReached, nil)
	assert.Equal(t, expected, actual)
}

func TestGoesPastStop(t *testing.T) {
	stopDistances := map[stopdistance.Key]float64{
		{"M55", 0, "1", "2"}: 250.0,
	}
	goesPastStop := map[bus.DirectedRoute][]bus.VehicleJourney{
		{"M55", 0, "ABC"}: {
			{
				DistanceFromStop: null.IntFrom(200), StopPointRef: null.StringFrom("1"),
				Timestamp: nulltypes.TimestampFrom(database.Timestamp{Time: time.Date(2019, 04, 23, 16, 38, 00, 0, time.UTC)}),
			},
			{
				DistanceFromStop: null.IntFrom(150), StopPointRef: null.StringFrom("1"),
				Timestamp: nulltypes.TimestampFrom(database.Timestamp{Time: time.Date(2019, 04, 23, 16, 39, 00, 0, time.UTC)}),
			},
			{
				DistanceFromStop: null.IntFrom(0), StopPointRef: null.StringFrom("1"),
				Timestamp: nulltypes.TimestampFrom(database.Timestamp{Time: time.Date(2019, 04, 23, 16, 40, 00, 0, time.UTC)}),
			},
			{
				DistanceFromStop: null.IntFrom(0), StopPointRef: null.StringFrom("2"),
				Timestamp: nulltypes.TimestampFrom(database.Timestamp{Time: time.Date(2019, 04, 23, 16, 50, 00, 0, time.UTC)}),
			},
		},
	}
	expected := getExpectedLabelledJourneys(goesPastStop, []int{120, 60, 600})
	actual := labels.Create(goesPastStop, stopDistances)
	assert.Equal(t, expected, actual)
}

func getExpectedLabelledJourneys(sequence map[bus.DirectedRoute][]bus.VehicleJourney, labels []int) (expected []bus.LabelledJourney) {
	for _, sequence := range sequence {
		for i, label := range labels {
			expected = append(expected, bus.LabelledJourneyFrom(sequence[i], label))
		}
	}
	return expected
}
