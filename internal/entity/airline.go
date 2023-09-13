package entity

import (
	"database/sql"
	"database/sql/driver"
	"encoding/json"
	"errors"
	"regexp"

	"github.com/metafates/smartway-test/internal/pkg/hashset"
)

var (
	_ sql.Scanner   = (*AirlineCode)(nil)
	_ driver.Valuer = (*AirlineCode)(nil)
)

var airlineCodePattern = regexp.MustCompile(`^[0-9A-Z\x{0400}-\x{04FF}]{2}$`)

type AirlineCode string

func (a *AirlineCode) Value() (driver.Value, error) {
	return string(*a), nil
}

func (a *AirlineCode) Scan(src any) error {
	if src == nil {
		*a = ""
		return nil
	}

	sv, err := driver.String.ConvertValue(src)
	if err != nil {
		return err
	}

	value, ok := sv.(string)
	if !ok {
		return errors.New("failed to scan airline id")
	}

	*a = AirlineCode(value)
	return nil
}

func (a *AirlineCode) UnmarshalText(text []byte) error {
	matches := airlineCodePattern.Match(text)

	if !matches {
		return errors.New("airline code must be two [A-ZА-Я0-9]")
	}

	*a = AirlineCode(text)
	return nil
}

func (a *AirlineCode) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}

	return a.UnmarshalText([]byte(s))
}

type Airline struct {
	Code AirlineCode `json:"code,omitempty"`
	Name string      `json:"name,omitempty"`
}

// AirlineChanges that can be applied
type AirlineChanges struct {
	Providers *hashset.Set[ProviderID] `json:"providers,omitempty"`
}
