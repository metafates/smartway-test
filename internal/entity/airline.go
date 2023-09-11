package entity

type Airline struct {
	Code string `json:"code,omitempty"`
	Name string `json:"name,omitempty"`

	// ProvidersIDs that provide this airline
	ProvidersIDs map[string]struct{} `json:"providers"`
}
