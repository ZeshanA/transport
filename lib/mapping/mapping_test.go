package mapping_test

import (
	"testing"
	"transport/lib/mapping"
	"transport/lib/testhelper"

	"github.com/stretchr/testify/assert"
)

func TestClient_RoadDistance(t *testing.T) {
	ts := testhelper.ServeMock(`{"code":"Ok","distances":[[1234.5]]}`)
	defer ts.Close()
	mc := mapping.NewClient("TEST", mapping.CustomBaseURLOption(ts.URL))
	expected := 1234.5
	assert.Equal(t, expected, mc.RoadDistance(1.2, 3.4, 5.6, 7.8))
}

func TestClient_RoadDistanceInvalidResponse(t *testing.T) {
	ts := testhelper.ServeMock(`{"code":"InvalidInput","distances":[[1234.5]]}`)
	defer ts.Close()
	errMsg := "mapping.RoadDistance: invalidInput response from MapBox API didn't cause the expected panic"
	mc := mapping.NewClient("TEST", mapping.CustomBaseURLOption(ts.URL))
	assert.Panics(t, func() { mc.RoadDistance(1.2, 3.4, 5.6, 7.8) }, errMsg)
}
