package mapping_test

import (
	"fmt"
	"testing"
	"transport/lib/mapping"
	"transport/lib/testhelper"

	"googlemaps.github.io/maps"

	"github.com/stretchr/testify/assert"
)

func TestRoadDistance(t *testing.T) {
	ts := testhelper.ServeMock(`{"rows": [{"elements": [{"distance": {"value": 1234}}]}], "status": "OK"}`)
	defer ts.Close()
	mc, err := maps.NewClient(maps.WithAPIKey("TEST"), maps.WithBaseURL(ts.URL))
	if err != nil {
		assert.Fail(t, fmt.Sprintf("failed to initialise maps client: %s", err))
	}
	expected := 1234.00
	assert.Equal(t, expected, mapping.RoadDistance(mc, 1.2, 3.4, 5.6, 7.8))
}
