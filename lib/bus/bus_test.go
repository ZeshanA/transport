package bus

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPartitionJourneys(t *testing.T) {
	expected := map[DirectedRoute][]VehicleJourney{
		{ExampleVJs[0].LineRef.String, int(ExampleVJs[0].DirectionRef.Int64), ExampleVJs[0].VehicleRef.String}: {ExampleVJs[0]},
		{ExampleVJs[1].LineRef.String, int(ExampleVJs[1].DirectionRef.Int64), ExampleVJs[1].VehicleRef.String}: {ExampleVJs[1]},
	}
	actual := PartitionJourneys([]VehicleJourney{ExampleVJs[0], ExampleVJs[1]})
	assert.Equal(t, expected, actual)
}
