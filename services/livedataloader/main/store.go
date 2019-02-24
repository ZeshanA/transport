package main

import (
	"database/sql"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"
	"transport/lib/database"
	"transport/lib/stringhelper"

	"github.com/lib/pq"

	"github.com/tidwall/gjson"
)

const timeFormat = "2006-01-02 15:04:05"

func store(liveVehicleData *string, dataIncoming chan bool) {
	db := database.OpenDBConnection()
	for {
		<-dataIncoming
		parseAndStore(liveVehicleData, db)
		log.Println("Finished sending vehicle entries to DB")
	}
}

func parseAndStore(liveVehicleData *string, db *sql.DB) {
	// Extract vehicle activity from JSON string
	vehicleActivity := gjson.Get(*liveVehicleData, VehicleActivityPath)

	arr := vehicleActivity.Array()
	fmt.Printf("This many items: %d\n", len(arr))

	txn, err := db.Begin()
	if err != nil {
		log.Fatal(err)
	}

	stmt, err := txn.Prepare(pq.CopyIn(
		"arrivals",
		"latitude", "longitude", "timestamp", "vehicle_id",
		"distance_along", "direction_id", "phase", "route_id",
		"trip_id", "next_stop_distance", "next_stop_id",
	))

	if err != nil {
		log.Fatal(err)
	}

	// Construct a DB row from each vehicle activity entry and insert the row into the DB
	vehicleActivity.ForEach(func(_, activityEntry gjson.Result) bool {
		fieldValues, err := getFieldValues(&activityEntry)
		if err != nil {
			fmt.Printf("Error whilst parsing field values: %v\n", fieldValues)
			return true
		}
		fmt.Printf("%v\n", fieldValues)
		// Insert fields into DB
		_, err = stmt.Exec(stringhelper.SliceToInterface(fieldValues)...)
		if err != nil {
			log.Printf("error occurred whilst executing insert statement for %v:\n%v\n", fieldValues, err)
			return true
		}
		return true
	})

	_, err = stmt.Exec()
	if err != nil {
		log.Printf("error whilst flushing statements: %v\n", err)
	}
	err = stmt.Close()
	if err != nil {
		log.Printf("error whilst closing insertion statement: %v\n", err)
	}
	err = txn.Commit()
	if err != nil {
		log.Printf("error whilst committing insertion transaction to db: %v\n", err)
	}

}

// Returns a slice representing a DB row that the given activityEntry corresponds to
func getFieldValues(activityEntry *gjson.Result) (*[]string, error) {
	// Some fields are nested directly under journey, others are nested more deeply
	journey := activityEntry.Get("MonitoredVehicleJourney")
	location := journey.Get("VehicleLocation")
	call := journey.Get("MonitoredCall")

	fieldCount := 11
	fields := make([]string, fieldCount)

	// Store fields in slice
	fields[0], fields[1] = numericNull(location.Get("Latitude").String()), numericNull(location.Get("Longitude").String())
	timestamp, err := parseTime(activityEntry.Get("RecordedAtTime").String())
	if err != nil {
		return nil, err
	}
	fields[2] = timestamp
	vehicleID, err := intFromID(journey.Get("VehicleRef").String())
	if err != nil {
		return nil, err
	}
	fields[3] = numericNull(strconv.Itoa(vehicleID))
	fields[4] = "0"
	fields[5] = numericNull(journey.Get("DirectionRef").String())
	phase, err := parseProgressRate(journey.Get("ProgressRate").String())
	if err != nil {
		return nil, err
	}
	fields[6] = phase
	fields[7] = journey.Get("LineRef").String()
	fields[8] = journey.Get("BlockRef").String()
	fields[9] = numericNull(call.Get("DistanceFromStop").String())
	fields[10] = numericNull(call.Get("StopPointRef").String())
	return &fields, nil
}

func numericNull(val string) string {
	if val == "" {
		return "0"
	}
	return val
}

func intFromID(stringID string) (numericID int, parseError error) {
	split := strings.Split(stringID, "_")
	numericValue, err := strconv.Atoi(split[1])
	if err != nil {
		return 0, err
	}
	return numericValue, nil
}

/*
Indicator of whether the bus is:
	- making progress (normalProgress) (i.e. moving, generally),
	- not moving (with value noProgress),
	- laying over before beginning a trip (value layover),
	- or serving a trip prior to one which will arrive (prevTrip).
*/
func parseProgressRate(progressRate string) (string, error) {
	switch progressRate {
	case "normalProgress", "noProgress":
		return "IN_PROGRESS", nil
	case "layover":
		return "LAYOVER_DURING", nil
	case "prevTrip":
		return "PREV_TRIP", nil
	default:
		return "", fmt.Errorf("invalid progress rate: %s", progressRate)
	}
}

func parseTime(timeString string) (string, error) {
	parsed, err := time.Parse(time.RFC3339, timeString)
	if err != nil {
		return "", fmt.Errorf("invalid timestamp: %s", timeString)
	}
	return parsed.Format(timeFormat), nil
}
