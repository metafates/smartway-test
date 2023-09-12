package entity

import (
	"encoding/json"
	"errors"
	"strconv"
)

type AccountID int

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

	// Schema is a Schema that this account is assigned to
	Schema SchemaID `json:"schemaId"`
}
