package models

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
)

// StringSlice is a custom type for []string to be stored as JSON
type StringSlice []string

// Value converts StringSlice to a database value (JSON)
func (s StringSlice) Value() (driver.Value, error) {
	if len(s) == 0 {
		return "[]", nil
	}
	return json.Marshal(s)
}

// Scan converts a database value (JSON) to StringSlice
func (s *StringSlice) Scan(value interface{}) error {
	bytes, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}
	return json.Unmarshal(bytes, s)
}
