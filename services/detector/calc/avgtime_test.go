package calc

import (
	"testing"
	"transport/lib/bus"

	"github.com/stretchr/testify/assert"
	"gopkg.in/guregu/null.v3"
)

func TestSplitMovementsByVehicleID(t *testing.T) {
	mvmts := []bus.LabelledJourney{
		{VehicleRef: null.StringFrom("A"), StopPointRef: null.StringFrom("A1")},
		{VehicleRef: null.StringFrom("B"), StopPointRef: null.StringFrom("B1")},
		{VehicleRef: null.StringFrom("B"), StopPointRef: null.StringFrom("B2")},
	}
	expected := map[string][]bus.LabelledJourney{
		"A": {
			{VehicleRef: null.StringFrom("A"), StopPointRef: null.StringFrom("A1")},
		},
		"B": {
			{VehicleRef: null.StringFrom("B"), StopPointRef: null.StringFrom("B1")},
			{VehicleRef: null.StringFrom("B"), StopPointRef: null.StringFrom("B2")},
		},
	}
	actual := SplitMovementsByVehicleID(mvmts)
	assert.Equal(t, expected, actual)
}
