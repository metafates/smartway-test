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
	_ sql.Scanner   = (*ProviderID)(nil)
	_ driver.Valuer = (*ProviderID)(nil)
)

type ProviderID string

func (p *ProviderID) Value() (driver.Value, error) {
	return string(*p), nil
}

func (p *ProviderID) Scan(src any) error {
	if src == nil {
		*p = ""
		return nil
	}

	sv, err := driver.String.ConvertValue(src)
	if err != nil {
		return err
	}

	value, ok := sv.(string)
	if !ok {
		return errors.New("failed to scan provider id")
	}

	*p = ProviderID(value)
	return nil
}

var providerIDPattern = regexp.MustCompile(`^[A-Z]{2}$`)

func (p *ProviderID) UnmarshalText(text []byte) error {
	id := string(text)

	matches := providerIDPattern.Match(text)
	if !matches {
		return errors.New("provider id must be two A-Z symbols")
	}

	*p = ProviderID(id)
	return nil
}

func (p *ProviderID) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}

	return p.UnmarshalText([]byte(s))
}

type Provider struct {
	ID   ProviderID `json:"id,omitempty"`
	Name string     `json:"name,omitempty"`

	// Airlines that this provider provides
	Airlines *hashset.Set[AirlineCode] `json:"-"`
}

// ProviderChanges that can be applied
type ProviderChanges struct {
	Name *string `json:"name,omitempty"`
}
