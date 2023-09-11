package entity

type Schema struct {
	Name      string              `json:"name,omitempty"`
	ID        string              `json:"id"`
	Providers map[string]struct{} `json:"suppliers"`
}
