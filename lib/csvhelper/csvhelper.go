package csvhelper

import (
	"fmt"
	"io/ioutil"
	"strings"
)

// RemoveNullRows deletes any rows with a "NULL" column value from
// the given CSV file and returns a byte[] containing only the valid rows
func RemoveNullRows(path string, columnSeparator string) (rows *[]byte, e error) {

	// Load file into memory
	fmt.Printf("Loading %s into memory to remove null rows...\n", path)
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("csv.RemoveNullRows: opening file '%s' failed due to: %v\n", path, err)
	}

	// Remove any containing null values, in-memory
	_, validRows := getValidRows(&data, columnSeparator)

	// Join all the valid rows into a single byte array, separated by new lines
	joinedRows := []byte(strings.Join(validRows, "\n"))

	return &joinedRows, nil
}

// getValidRows returns a new slice containing all the CSV rows in 'data' that
// contained no NULL values in any column
func getValidRows(data *[]byte, columnSeparator string) (nullRowCount int, validRowsSlice []string) {
	rows := strings.Split(string(*data), "\n")
	nullRows := 0
	var validRows []string

	for i, row := range rows {
		fmt.Printf("Removing null rows: processed %d rows of %d...\n", i+1, len(rows))
		validRow := true
		columns := strings.Split(row, columnSeparator)
		for _, col := range columns {
			if col == "NULL" {
				validRow = false
				nullRows++
				break
			}
		}
		if validRow {
			validRows = append(validRows, row)
		}
	}
	fmt.Printf("Succesfully removed %d null rows...\n", nullRows)
	return nullRows, validRows
}
