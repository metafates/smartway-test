package entity

import (
	"database/sql"
	"database/sql/driver"
	"encoding/json"
	"errors"
	"strconv"

	"github.com/metafates/smartway-test/internal/pkg/hashset"
)

var (
	_ sql.Scanner   = (*SchemaID)(nil)
	_ driver.Valuer = (*SchemaID)(nil)
)

type SchemaID int

func (s *SchemaID) Value() (driver.Value, error) {
	return int64(*s), nil
}

func (s *SchemaID) Scan(src any) error {
	if src == nil {
		*s = 0
		return nil
	}

	iv, err := driver.Int32.ConvertValue(src)
	if err != nil {
		return err
	}

	value, ok := iv.(int)
	if !ok {
		return errors.New("failed to scan schema id")
	}

	*s = SchemaID(value)
	return nil
}

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
}

// SchemaChanges that can be applied
type SchemaChanges struct {
	Name      *string                  `json:"name,omitempty"`
	Providers *hashset.Set[ProviderID] `json:"providers,omitempty"`
}
