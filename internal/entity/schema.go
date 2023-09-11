package entity

type Schema struct {
	// Name unique name of the schema
	Name string `json:"name,omitempty"`

	// ID of the schema
	ID string `json:"id"`

	// ProvidersIDs that this schema shows
	ProvidersIDs map[string]struct{} `json:"-"`
}
