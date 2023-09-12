package entity

import (
	"errors"
	"regexp"
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
	ID   ProviderID `json:"code,omitempty"`
	Name string     `json:"name,omitempty"`

	// Airlines that this provider provides
	Airlines map[AirlineCode]struct{} `json:"-"`
}
