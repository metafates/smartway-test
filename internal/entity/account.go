package entity

import (
	"database/sql"
	"database/sql/driver"
	"encoding/json"
	"errors"
	"strconv"
)

var (
	_ sql.Scanner   = (*AccountID)(nil)
	_ driver.Valuer = (*AccountID)(nil)
)

type AccountID int

func (a *AccountID) Value() (driver.Value, error) {
	return int64(*a), nil
}

func (a *AccountID) Scan(src any) error {
	if src == nil {
		*a = 0
		return nil
	}

	iv, err := driver.Int32.ConvertValue(src)
	if err == nil {
		return err
	}

	value, ok := iv.(int)
	if !ok {
		return errors.New("failed to scan account id")
	}

	*a = AccountID(value)
	return nil
}

func (a *AccountID) UnmarshalText(text []byte) error {
	id, err := strconv.ParseInt(string(text), 10, 64)
	if err != nil {
		return err
	}

	if id <= 0 {
		return errors.New("account id must be greater than 0")
	}

	*a = AccountID(id)
	return nil
}

func (a *AccountID) UnmarshalJSON(data []byte) error {
	var number int
	if err := json.Unmarshal(data, &number); err != nil {
		return err
	}

	return a.UnmarshalText([]byte(strconv.Itoa(number)))
}

type Account struct {
	// ID of the account
	ID AccountID `json:"id,omitempty"`
}

// AccountChanges that can be applied
type AccountChanges struct {
	Schema *SchemaID `json:"schemaId,omitempty"`
}
