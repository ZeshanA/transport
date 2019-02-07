package csvhelper

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"
)

// RemoveNullRows deletes any rows with a "NULL" column value from
// the given CSV file
func RemoveNullRows(path string, columnSeparator string) (nullRowCount int, e error) {

	// Load file into memory
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return 0, fmt.Errorf("csv.RemoveNullRows: opening file '%s' failed due to: %v\n", path, err)
	}

	// Remove any in-memory rows containing null values
	nullRows, validRows := getValidRows(&data, columnSeparator)

	// Empty the file so we can overwrite all rows with the cleaned data in memory
	fmt.Printf("Emptying file '%s'...\n", path)
	err = os.Truncate(path, 0)
	if err != nil {
		return 0, fmt.Errorf("csv.RemoveNullRows: emptying file '%s' failed due to: %v\n", path, err)
	}
	fmt.Printf("Succesfully emptied file '%s'...\n", path)

	// Write just the valid, non-null rows back to the file
	fmt.Printf("Writing %d valid rows to file '%s'...\n", len(validRows), path)
	err = ioutil.WriteFile(path, []byte(strings.Join(validRows, "\n")), 0666)
	if err != nil {
		return 0, fmt.Errorf("csv.RemoveNullRows: writing valid rows to file '%s' failed due to: %v\n", path, err)
	}

	return nullRows, nil
}

func getValidRows(data *[]byte, columnSeparator string) (nullRowCount int, validRowsSlice []string) {
	rows := strings.Split(string(*data), "\n")
	nullRows := 0
	var validRows []string

	for i, row := range rows {
		fmt.Printf("Removing null rows from %d of %d...\n", i+1, len(rows))
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
	return nullRows, validRows
}
