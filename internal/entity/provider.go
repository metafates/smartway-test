package entity

import (
	"encoding/json"
	"errors"
	"regexp"

	"github.com/metafates/smartway-test/internal/pkg/hashset"
)

type ProviderID string

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
