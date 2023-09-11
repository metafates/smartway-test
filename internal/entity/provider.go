package entity

type Provider struct {
	ID   string `json:"code,omitempty"`
	Name string `json:"name,omitempty"`

	// AirlinesCodes that this provider provides
	AirlinesCodes map[string]struct{} `json:"airlines"`
}
