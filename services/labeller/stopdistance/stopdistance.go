package stopdistance

import (
	"database/sql"
	"log"
	"transport/lib/bus"
	"transport/lib/database"
)

// stopdistance.Get fetches all stop distances from
// the database and returns a map of distances,
// queryable by StopDistanceKey.
func Get(db *sql.DB) map[Key]float64 {
	rows := database.FetchAllRows(db, database.StopDistanceTable.Name)
	defer rows.Close()
	sdList := scanIntoStructs(rows)
	return partitionStopDistanceList(sdList)
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

func GetAverage(db *sql.DB) map[string]int {
	// Init map
	distances := map[string]int{}
	// Fetch rows from average distance table
	rows := database.FetchAllRows(db, database.AverageDistanceTable.Name)
	defer rows.Close()

	for rows.Next() {
		// Scan row into relevant variables
		var routeID string
		var distance int
		err := rows.Scan(&routeID, &distance)
		if err != nil {
			log.Fatal(err)
		}
		// Store distance under correct routeID in map
		distances[routeID] = distance
	}
	err := rows.Err()
	if err != nil {
		log.Fatal(err)
	}
	return distances
}
