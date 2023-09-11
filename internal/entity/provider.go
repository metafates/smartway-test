package entity

type Provider struct {
	ID       string              `json:"code,omitempty"`
	Name     string              `json:"name,omitempty"`
	Airlines map[string]struct{} `json:"airlines"`
}
