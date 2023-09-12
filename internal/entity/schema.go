package entity

import (
	"errors"
	"strconv"
)

type SchemaID int

func (s *SchemaID) UnmarshalText(text []byte) error {
	id, err := strconv.ParseInt(string(text), 10, 64)
	if err != nil {
		return err
	}

	if id <= 0 {
		return errors.New("schema id must be greater than 0")
	}

	*s = SchemaID(id)
	return nil
}

func (s *SchemaID) UnmarshalJSON(data []byte) error {
	return s.UnmarshalText(data)
}

type Schema struct {
	// Name unique name of the schema
	Name string `json:"name,omitempty"`

	// ID of the schema
	ID SchemaID `json:"id"`

	// ProvidersIDs that this schema shows
	ProvidersIDs map[ProviderID]struct{} `json:"-"`
}
