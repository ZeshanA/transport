package nulltypes

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"strings"
	"time"
	"transport/lib/database"
)

// null.Timestamp
type Timestamp struct {
	database.Timestamp
	Valid bool
}

func TimestampFrom(ts database.Timestamp) Timestamp {
	if ts.IsZero() {
		return Timestamp{ts, false}
	}
	return Timestamp{ts, true}
}

// MarshalJSON converts a null.Timestamp into a JSON []byte
func (ts *Timestamp) MarshalJSON() ([]byte, error) {
	if ts.Valid {
		return []byte(fmt.Sprintf(`"%s"`, ts.Timestamp.Format(database.TimeFormat))), nil
	}
	return []byte("null"), nil
}

// UnarshalJSON converts a JSON []byte into a null.Timestamp
func (ts *Timestamp) UnmarshalJSON(b []byte) error {
	// Remove quote marks
	str := strings.Replace(string(b), `"`, "", -1)
	if str == "null" {
		ts.Valid = false
		return nil
	}
	t, err := time.Parse(database.TimeFormat, str)
	if err != nil {
		return err
	}
	ts.Timestamp = database.Timestamp{Time: t}
	return nil
}

// Scan uses a cell from the DB to populate a Timestamp struct
func (ts *Timestamp) Scan(value interface{}) error {
	if value == nil {
		ts.Timestamp, ts.Valid = database.Timestamp{}, false
		return nil
	}
	switch v := value.(type) {
	case string:
		if v == "NULL" {
			ts.Timestamp, ts.Valid = database.Timestamp{}, false
			return nil
		}
		parsed, err := time.Parse(database.TimeFormat, v)
		if err != nil {
			return fmt.Errorf("nulltypes.Timestamp.Scan: unable to parse string: %v\n", value)
		}
		ts.Timestamp, ts.Valid = database.Timestamp{Time: parsed}, true
	case time.Time:
		ts.Timestamp, ts.Valid = database.Timestamp{Time: v}, true
	default:
		return fmt.Errorf("nulltypes.Timestamp.Scan: invalid type passed in: %v\n", v)
	}
	return nil
}

// Value takes a Timestamp struct and outputs a value that can be stored
// in the DB
func (ts Timestamp) Value() (driver.Value, error) {
	if !ts.Valid {
		return nil, nil
	}
	return ts.Timestamp.Time, nil
}

// null.StringSlice
type StringSlice struct {
	StringSlice []string
	Valid       bool
}

func StringSliceFrom(ss []string) StringSlice {
	return StringSlice{ss, true}
}

func (ss *StringSlice) MarshalJSON() ([]byte, error) {
	if ss.Valid {
		return json.Marshal(ss.StringSlice)
	}
	return []byte(`[]`), nil
}

func (ss *StringSlice) UnmarshalJSON(b []byte) error {
	if string(b) == "null" {
		ss.Valid = false
		return nil
	}
	// Unmarshal into a []string
	var arr []string
	err := json.Unmarshal(b, &arr)
	if err != nil {
		return err
	}
	// Store []string inside struct and mark valid
	ss.StringSlice, ss.Valid = arr, true
	return nil
}

func (ss *StringSlice) Scan(value interface{}) error {
	if value == nil {
		ss.StringSlice, ss.Valid = nil, false
		return nil
	}
	ss.Valid = true
	switch v := value.(type) {
	case string:
		ss.StringSlice = parseStringSlice([]byte(v))
	case []byte:
		ss.StringSlice = parseStringSlice(v)
	default:
		ss.StringSlice, ss.Valid = nil, false
		return fmt.Errorf("nulltypes.StringSlice.Scan: invalid type passed in: %v\n", value)
	}
	return nil
}

func parseStringSlice(rawBytes []byte) []string {
	str := string(rawBytes)
	// Remove {} braces from start and end of string
	noBrackets := str[1 : len(str)-1]
	// Split the string into its components
	components := strings.Split(noBrackets, ",")
	// Remove quote marks from each component
	for i, v := range components {
		components[i] = strings.Replace(v, `"`, ``, -1)
	}
	return components
}

func (ss *StringSlice) Value() (driver.Value, error) {
	if !ss.Valid {
		return nil, nil
	}
	return ss.StringSlice, nil
}
