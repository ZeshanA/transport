package stopdistance

import (
	"database/sql"
	"fmt"
	"log"
	"transport/lib/bus"
	"transport/lib/database"
)

// Get fetches all stop distances from
// the database and returns a map of distances,
// queryable by StopDistanceKey.
func Get(db *sql.DB) map[Key]float64 {
	rows := fetchRows(db)
	defer rows.Close()
	sdList := scanIntoStructs(rows)
	return partitionStopDistanceList(sdList)
}

// Fetch raw rows from stop_distance table
func fetchRows(db *sql.DB) *sql.Rows {
	rows, err := db.Query(fmt.Sprintf(`SELECT * FROM %s`, database.StopDistanceTable.Name))
	if err != nil {
		log.Panicf("stopdistance.Get: error whilst reading from db: %s\n", err)
	}
	return rows
}

// Convert database rows into bus.StopDistance structs
func scanIntoStructs(rows *sql.Rows) []bus.StopDistance {
	var sdList []bus.StopDistance
	for rows.Next() {
		var sd bus.StopDistance
		err := rows.Scan(&sd.RouteID, &sd.FromID, &sd.ToID, &sd.Distance, &sd.DirectionID)
		if err != nil {
			log.Fatal(err)
		}
		sdList = append(sdList, sd)
	}
	err := rows.Err()
	if err != nil {
		log.Fatal(err)
	}
	return sdList
}

type Key struct {
	RouteID     string
	DirectionID int
	FromID      string
	ToID        string
}

// Returns a nested map of fromStopID -> toStopID -> distance
func partitionStopDistanceList(distances []bus.StopDistance) map[Key]float64 {
	partitioned := map[Key]float64{}
	for _, sd := range distances {
		key := Key{RouteID: sd.RouteID, DirectionID: sd.DirectionID, FromID: sd.FromID, ToID: sd.ToID}
		partitioned[key] = sd.Distance
	}
	return partitioned
}
