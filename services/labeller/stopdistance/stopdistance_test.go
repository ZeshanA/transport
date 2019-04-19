package stopdistance

import (
	"testing"
	"transport/lib/bus"

	"github.com/stretchr/testify/assert"

	"github.com/DATA-DOG/go-sqlmock"
)

func TestShouldGetStopDistances(t *testing.T) {
	// Set up DB mock
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	// Rows for mock DB to return
	rows := sqlmock.NewRows([]string{"route_id", "from_id", "to_id", "distance", "direction_id"}).
		AddRow("route1", "stop1", "stop2", 123, 0).
		AddRow("route1", "stop2", "stop3", 456, 0)

	// Expect the correct query to be executed
	mock.ExpectQuery(".\\*").WillReturnRows(rows)

	// Verify the final returned slice of structs is as expected
	expected := []bus.StopDistance{
		{"route1", 0, "stop1", "stop2", 123},
		{"route1", 0, "stop2", "stop3", 456},
	}
	actual := Get(db)
	assert.Equal(t, expected, actual)

	// Verify the correct query was executed
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("TestShouldGetStopDistances: there were unfulfilled expectations: %s", err)
	}
}
