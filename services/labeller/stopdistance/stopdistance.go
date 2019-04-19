package stopdistance

import (
	"database/sql"
	"fmt"
	"log"
	"transport/lib/bus"
	"transport/lib/database"
)

// Get fetches all stop distances from
// the database and returns a slice of them, parsed
// into bus.StopDistance structs.
func Get(db *sql.DB) []bus.StopDistance {
	rows := fetchRows(db)
	defer rows.Close()
	sdList := scanIntoStructs(rows)
	return sdList
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
