package entity

type Schema struct {
	// Name unique name of the schema
	Name string `json:"name,omitempty"`

	// ID of the schema
	ID int `json:"id"`

	// ProvidersIDs that this schema shows
	ProvidersIDs map[int]struct{} `json:"-"`
}
