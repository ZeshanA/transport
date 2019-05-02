package stopdistance

import (
	"database/sql/driver"
	"testing"
	"transport/lib/database"
	"transport/lib/testhelper"

	"github.com/stretchr/testify/assert"
)

func TestShouldGetStopDistances(t *testing.T) {
	var mockRows = [][]driver.Value{
		{"route1", "stop1", "stop2", 123.0, 0},
		{"route1", "stop2", "stop3", 456.0, 0},
	}
	db, mock := testhelper.SetupDBMock(t, database.StopDistanceTable.Columns, mockRows, ".\\*")
	defer db.Close()

	// Verify the final returned slice of structs is as expected
	expectedStructs := map[Key]float64{
		{"route1", 0, "stop1", "stop2"}: 123.0,
		{"route1", 0, "stop2", "stop3"}: 456.0,
	}
	actual := Get(db)
	assert.Equal(t, expectedStructs, actual)

	// Verify the correct query was executed
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("TestShouldGetStopDistances: there were unfulfilled expectations: %s", err)
	}
}

func TestShouldGetAverageStopDistances(t *testing.T) {
	var mockRows = [][]driver.Value{
		{"route1", 123},
		{"route2", 456},
	}
	db, mock := testhelper.SetupDBMock(t, database.AverageDistanceTable.Columns, mockRows, ".\\*")
	defer db.Close()

	// Verify the final returned slice of structs is as expected
	expectedStructs := map[string]int{
		"route1": 123,
		"route2": 456,
	}
	actual := GetAverage(db)
	assert.Equal(t, expectedStructs, actual)

	// Verify the correct query was executed
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("TestShouldGetAverageStopDistances: there were unfulfilled expectations: %s", err)
	}
}
