package usecase

import (
	"context"
	"errors"

	"github.com/metafates/smartway-test/internal/entity"
	"github.com/samber/lo"
)

var (
	ErrAccountNotFound = errors.New("account not found")
)

var _ Account = (*AccountUseCase)(nil)

type AccountUseCase struct {
	repo Repository
}

func NewAccountUseCase(repository Repository) *AccountUseCase {
	return &AccountUseCase{repo: repository}
}

func (a *AccountUseCase) Add(ctx context.Context, account entity.Account) error {
	return a.repo.StoreAccount(ctx, account)
}

func (a *AccountUseCase) Delete(ctx context.Context, ID string) error {
	return a.repo.DeleteAccount(ctx, ID)
}

func (a *AccountUseCase) SetSchema(ctx context.Context, accountID, schemaID string) error {
	account, accountExists, err := a.repo.GetAccountByID(ctx, accountID)
	if err != nil {
		return err
	}

	if !accountExists {
		return ErrAccountNotFound
	}

	_, schemaFound, err := a.repo.GetSchemaByID(ctx, schemaID)
	if err != nil {
		return err
	}

	if !schemaFound {
		return ErrSchemaNotFound
	}

	account.SchemaID = schemaID

	return a.repo.StoreAccount(ctx, account)
}

func (a *AccountUseCase) GetAirlines(ctx context.Context, ID string) ([]entity.Airline, error) {
	account, accountExists, err := a.repo.GetAccountByID(ctx, ID)
	if err != nil {
		return nil, err
	}

	if !accountExists {
		return nil, ErrAccountNotFound
	}

	schema, schemaFound, err := a.repo.GetSchemaByID(ctx, account.SchemaID)
	if err != nil {
		return nil, err
	}

	if !schemaFound {
		return nil, ErrSchemaNotFound
	}

	providersIDs := lo.Keys(schema.ProvidersIDs)

	providers, err := a.repo.GetProvidersByIDs(ctx, providersIDs...)
	if err != nil {
		return nil, err
	}

	var airlinesCodes map[string]struct{}
	for _, provider := range providers {
		for code := range provider.AirlinesCodes {
			airlinesCodes[code] = struct{}{}
		}
	}

	return a.repo.GetAirlinesByCodes(ctx, lo.Keys(airlinesCodes)...)
}
