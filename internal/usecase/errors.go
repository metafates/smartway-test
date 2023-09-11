package usecase

import "errors"

var (
	ErrSchemaDoesNotExist  = errors.New("schema does not exist")
	ErrAccountDoesNotExist = errors.New("account does not exist")
)
