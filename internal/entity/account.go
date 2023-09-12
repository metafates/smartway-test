package entity

type Account struct {
	// ID of the account
	ID int `json:"id,omitempty"`

	// SchemaID is a SchemaID id what this account is assigned to
	SchemaID int `json:"-"`
}
