package entity

import (
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

type Account struct {
	// ID of the account
	ID AccountID `json:"id,omitempty"`

	// Schema is a Schema that this account is assigned to
	Schema SchemaID `json:"schemaId"`
}
