package usecase

import "errors"

var (
	ErrSchemaNotFound   = errors.New("schema not found")
	ErrAccountNotFound  = errors.New("account not found")
	ErrProviderNotFound = errors.New("provider not found")
)
