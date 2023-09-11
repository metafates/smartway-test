package usecase

import (
	"context"

	"github.com/metafates/smartway-test/internal/entity"
	"github.com/samber/lo"
)

var _ Account = (*AccountUseCase)(nil)

type AccountUseCase struct {
	repo Repository
}

func (a *AccountUseCase) Add(ctx context.Context, account entity.Account) error {
	return a.repo.StoreAccount(ctx, account, false)
}

func (a *AccountUseCase) Delete(ctx context.Context, ID string) error {
	return a.repo.DeleteAccount(ctx, ID)
}

func (a *AccountUseCase) SetScheme(ctx context.Context, accountID, schemaID string) error {
	account, accountExists, err := a.repo.GetAccount(ctx, accountID)
	if err != nil {
		return err
	}

	if !accountExists {
		return ErrAccountDoesNotExist
	}

	_, schemaExists, err := a.repo.GetSchema(ctx, schemaID)
	if err != nil {
		return err
	}

	if !schemaExists {
		return ErrSchemaDoesNotExist
	}

	account.SchemaID = schemaID

	return a.repo.StoreAccount(ctx, account, true)
}

func (a *AccountUseCase) GetAirlines(ctx context.Context, ID string) ([]entity.Airline, error) {
	account, accountExists, err := a.repo.GetAccount(ctx, ID)
	if err != nil {
		return nil, err
	}

	if !accountExists {
		return nil, ErrAccountDoesNotExist
	}

	schema, schemaExists, err := a.repo.GetSchema(ctx, account.SchemaID)
	if err != nil {
		return nil, err
	}

	if !schemaExists {
		return nil, ErrSchemaDoesNotExist
	}

	providersIDs := lo.Keys(schema.Providers)

	providers, err := a.repo.GetProviders(ctx, providersIDs...)
	if err != nil {
		return nil, err
	}

	var airlinesIDs map[string]struct{}
	for _, provider := range providers {
		for airlineID := range provider.Airlines {
			airlinesIDs[airlineID] = struct{}{}
		}
	}

	return a.repo.GetAirlines(ctx, lo.Keys(airlinesIDs)...)
}

func NewAccountUseCase(repository Repository) *AccountUseCase {
	return &AccountUseCase{repo: repository}
}
