package entity

import (
	"encoding/json"
	"errors"
	"strconv"

	"github.com/metafates/smartway-test/internal/pkg/hashset"
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
	var number int
	if err := json.Unmarshal(data, &number); err != nil {
		return err
	}

	return s.UnmarshalText([]byte(strconv.Itoa(number)))
}

type Schema struct {
	// Name unique name of the schema
	Name string `json:"name"`

	// ID of the schema
	ID SchemaID `json:"id"`

	// Providers that this schema shows
	Providers *hashset.Set[ProviderID] `json:"providers"`
}

// SchemaChanges that can be applied
type SchemaChanges struct {
	Name      *string                  `json:"name,omitempty"`
	Providers *hashset.Set[ProviderID] `json:"providers,omitempty"`
}
