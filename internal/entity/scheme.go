package entity

type Scheme struct {
	Name      string              `json:"name,omitempty"`
	Suppliers map[string]Supplier `json:"suppliers"`
}
