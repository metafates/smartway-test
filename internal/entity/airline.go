package entity

import (
	"encoding/json"
	"errors"
	"regexp"

	"github.com/metafates/smartway-test/internal/pkg/hashset"
)

type AirlineCode string

var airlineCodePattern = regexp.MustCompile(`^[0-9A-Z\x{0400}-\x{04FF}]{2}$`)

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

	// Providers that provide this airline
	Providers *hashset.Set[ProviderID] `json:"-"`
}
