package entity

type Supplier struct {
	Code     string             `json:"code,omitempty"`
	Name     string             `json:"name,omitempty"`
	Airlines map[string]Airline `json:"airlines"`
}
