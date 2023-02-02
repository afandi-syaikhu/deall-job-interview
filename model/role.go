package model

import (
	"database/sql/driver"
	"errors"
)

type Role string

func (r *Role) Scan(value interface{}) error {
	asBytes, ok := value.([]byte)
	if !ok {
		return errors.New("Scan source is not []byte")
	}
	*r = Role(string(asBytes))
	return nil
}

func (r Role) Value() (driver.Value, error) {
	// validation would go here
	return string(r), nil
}
