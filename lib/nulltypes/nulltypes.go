package nulltypes

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"log"
	"time"
	"transport/lib/database"
)

// null.Timestamp
type Timestamp struct {
	Timestamp database.Timestamp
	Valid     bool
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
	t, err := time.Parse(database.TimeFormat, string(b))
	if err != nil {
		return err
	}
	ts.Timestamp = database.Timestamp{Time: t}
	return nil
}

func (ts *Timestamp) Scan(value interface{}) error {
	if value == nil {
		ts.Timestamp, ts.Valid = database.Timestamp{}, false
		return nil
	}
	ts.Valid = true
	newTS := database.Timestamp{}
	err := newTS.UnmarshalJSON(value.([]byte))
	ts.Timestamp = newTS
	if err != nil {
		log.Fatal(err)
	}
	return nil
}

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
	// Use JSON unmarshaller for convenience
	err := ss.UnmarshalJSON(value.([]byte))
	if err != nil {
		return err
	}
	return nil
}

func (ss *StringSlice) Value() (driver.Value, error) {
	if !ss.Valid {
		return nil, nil
	}
	return ss.StringSlice, nil
}
