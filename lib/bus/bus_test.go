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

func TestRemoveAgencyID(t *testing.T) {
	expected := map[string]string{
		"MTA NYCT_B41":   "B41",
		"MTA NYCT_BX41+": "BX41+",
		"MTABC_Q10":      "Q10",
		"MTABC_Q114":     "Q114",
	}
	for input, expectedOutput := range expected {
		actualOutput := RemoveAgencyID(input)
		assert.Equal(t, expectedOutput, actualOutput)
	}
}
