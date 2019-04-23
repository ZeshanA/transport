package fetch_test

import (
	"labeller/fetch"
	"testing"
	"time"
	"transport/lib/bus"
	"transport/lib/database"
	"transport/lib/testhelper"

	"github.com/stretchr/testify/assert"
)

var expectedQuery = `SELECT \* FROM ([^ ])* WHERE TIMESTAMP BETWEEN '([0-9]){4}-(([0-9]){1}|([0-9]){2})-(([0-9]){1}|([0-9]){2}) 00:00:00' AND '([0-9]){4}-(([0-9]){1}|([0-9]){2})-(([0-9]){1}|([0-9]){2}) 23:59:59'`

func TestShouldGetDateRange(t *testing.T) {
	colNames := database.VehicleJourneyTable.Columns
	colNames = append(colNames, "entry_id")
	db, mock := testhelper.SetupDBMock(t, colNames, bus.ExampleVJRows, expectedQuery)
	defer db.Close()

	// Verify the final returned slice of structs is as expected
	expected := [][]bus.VehicleJourney{{bus.ExampleVJs[0], bus.ExampleVJs[1]}}
	actual := fetch.DateRange(db, time.Now(), time.Now())
	assert.Equal(t, expected, actual)

	// Verify the correct query was executed
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("TestShouldGetStopDistances: there were unfulfilled expectations: %s", err)
	}
}
