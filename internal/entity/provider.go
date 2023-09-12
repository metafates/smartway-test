package entity

import (
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

type Provider struct {
	ID   ProviderID `json:"id,omitempty"`
	Name string     `json:"name,omitempty"`

	// Airlines that this provider provides
	Airlines *hashset.Set[AirlineCode] `json:"-"`
}
