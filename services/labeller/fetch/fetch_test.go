package fetch_test

import (
	"database/sql/driver"
	"labeller/fetch"
	"testing"
	"time"
	"transport/lib/bus"
	"transport/lib/database"
	"transport/lib/nulltypes"
	"transport/lib/testhelper"

	"gopkg.in/guregu/null.v3"

	"github.com/stretchr/testify/assert"
)

func TestShouldGetDateRange(t *testing.T) {
	colNames := database.VehicleJourneyTable.Columns
	colNames = append(colNames, "entry_id")
	db, mock := testhelper.SetupDBMock(t, colNames, mockRows, expectedQuery)
	defer db.Close()

	// Verify the final returned slice of structs is as expected
	actual := fetch.DateRange(db, time.Now(), time.Now())
	assert.Equal(t, expectedStructs, actual)

	// Verify the correct query was executed
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("TestShouldGetStopDistances: there were unfulfilled expectations: %s", err)
	}
}

var mockRows = [][]driver.Value{
	{
		"MTA NYCT_M55", 1, "MTA NYCT_MQ_E9-Sunday-038000_M55_2",
		"M55", "MTA NYCT", "MTA_803080", "MTA_803185", "NULL", `{"MTA NYCT_224082","MTA BC_220132"}`,
		40.721857, -73.999738, "normalProgress", "", "MTA NYCT_3814", "2019-04-21 06:37:47.238",
		"2019-04-21 06:37:47.238", 161, 0, "MTA_400159", "2019-04-21 06:37:13", 874970616,
	},
	{
		"MTA NYCT_M41", 0, "MTA NYCT_M41_E9-Sunday-038000_M41_2",
		"M41", "MTA NYCT", "MTA_123456", "MTA_67891", "NULL", `{"MTA NYCT_112233","MTA BC_445566"}`,
		39.91857, -58.319738, "normalProgress", "", "MTA NYCT_2418", "2019-03-22 05:51:31.338",
		"2019-03-22 05:51:39.338", 349, 0, "MTA_338991", "2019-03-22 05:31:31.338", 874970617,
	},
}

var expectedQuery = `SELECT \* FROM ([^ ])* WHERE TIMESTAMP BETWEEN '([0-9]){4}-(([0-9]){1}|([0-9]){2})-(([0-9]){1}|([0-9]){2}) 00:00:00' AND '([0-9]){4}-(([0-9]){1}|([0-9]){2})-(([0-9]){1}|([0-9]){2}) 23:59:59'`

var expectedStructs = [][]bus.VehicleJourney{
	{
		bus.VehicleJourney{
			LineRef:                  null.StringFrom("MTA NYCT_M55"),
			DirectionRef:             null.IntFrom(1),
			TripID:                   null.StringFrom("MTA NYCT_MQ_E9-Sunday-038000_M55_2"),
			PublishedLineName:        null.StringFrom("M55"),
			OperatorRef:              null.StringFrom("MTA NYCT"),
			OriginRef:                null.StringFrom("MTA_803080"),
			DestinationRef:           null.StringFrom("MTA_803185"),
			OriginAimedDepartureTime: nulltypes.Timestamp{},
			SituationRef:             nulltypes.StringSliceFrom([]string{"MTA NYCT_224082", "MTA BC_220132"}),
			Latitude:                 null.FloatFrom(-73.999738),
			Longitude:                null.FloatFrom(40.721857),
			ProgressRate:             null.StringFrom("normalProgress"),
			Occupancy:                null.StringFrom(""),
			VehicleRef:               null.StringFrom("MTA NYCT_3814"),
			ExpectedArrivalTime:      nulltypes.TimestampFrom(database.Timestamp{Time: time.Date(2019, 4, 21, 6, 37, 47, 238000000, time.UTC)}),
			ExpectedDepartureTime:    nulltypes.TimestampFrom(database.Timestamp{Time: time.Date(2019, 4, 21, 6, 37, 47, 238000000, time.UTC)}),
			DistanceFromStop:         null.IntFrom(161),
			NumberOfStopsAway:        null.IntFrom(0),
			StopPointRef:             null.StringFrom("MTA_400159"),
			Timestamp:                nulltypes.TimestampFrom(database.Timestamp{Time: time.Date(2019, 4, 21, 6, 37, 13, 0, time.UTC)}),
		},
		bus.VehicleJourney{
			LineRef:                  null.StringFrom("MTA NYCT_M41"),
			DirectionRef:             null.IntFrom(0),
			TripID:                   null.StringFrom("MTA NYCT_M41_E9-Sunday-038000_M41_2"),
			PublishedLineName:        null.StringFrom("M41"),
			OperatorRef:              null.StringFrom("MTA NYCT"),
			OriginRef:                null.StringFrom("MTA_123456"),
			DestinationRef:           null.StringFrom("MTA_67891"),
			OriginAimedDepartureTime: nulltypes.Timestamp{},
			SituationRef:             nulltypes.StringSliceFrom([]string{"MTA NYCT_112233", "MTA BC_445566"}),
			Latitude:                 null.FloatFrom(-58.319738),
			Longitude:                null.FloatFrom(39.91857),
			ProgressRate:             null.StringFrom("normalProgress"),
			Occupancy:                null.StringFrom(""),
			VehicleRef:               null.StringFrom("MTA NYCT_2418"),
			ExpectedArrivalTime:      nulltypes.TimestampFrom(database.Timestamp{Time: time.Date(2019, 3, 22, 5, 51, 31, 338000000, time.UTC)}),
			ExpectedDepartureTime:    nulltypes.TimestampFrom(database.Timestamp{Time: time.Date(2019, 3, 22, 5, 51, 39, 338000000, time.UTC)}),
			DistanceFromStop:         null.IntFrom(349),
			NumberOfStopsAway:        null.IntFrom(0),
			StopPointRef:             null.StringFrom("MTA_338991"),
			Timestamp:                nulltypes.TimestampFrom(database.Timestamp{Time: time.Date(2019, 3, 22, 5, 31, 31, 338000000, time.UTC)}),
		},
	},
}
