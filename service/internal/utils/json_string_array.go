package utils

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
)

// JSONStringArray is a custom type that wraps []string
type JSONStringArray []string

// Scan implements the sql.Scanner interface
func (a *JSONStringArray) Scan(value interface{}) error {
	if value == nil {
		*a = []string{}
		return nil
	}

	b, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}

	return json.Unmarshal(b, a)
}

// Value implements the driver.Valuer interface
func (a JSONStringArray) Value() (driver.Value, error) {
	return json.Marshal(a)
}
