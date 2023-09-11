package entity

type Account struct {
	// ID of the account
	ID string `json:"id,omitempty"`

	// SchemaID is a SchemaID id what this account is assigned to
	SchemaID string `json:"-"`
}
