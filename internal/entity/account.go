package entity

type Account struct {
	Name   string `json:"name,omitempty"`
	Scheme Scheme `json:"scheme"`
}
