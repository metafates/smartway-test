package usecase

import (
	"context"
	"errors"

	"github.com/metafates/smartway-test/internal/entity"
	"github.com/metafates/smartway-test/internal/pkg/hashset"
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

func (a *AccountUseCase) Delete(ctx context.Context, ID entity.AccountID) error {
	return a.repo.DeleteAccount(ctx, ID)
}

func (a *AccountUseCase) SetSchema(ctx context.Context, accountID entity.AccountID, schemaID entity.SchemaID) error {
	return a.repo.UpdateAccount(ctx, accountID, entity.AccountChanges{
		Schema: &schemaID,
	})
}

func (a *AccountUseCase) GetAirlines(ctx context.Context, ID entity.AccountID) ([]entity.Airline, error) {
	account, accountExists, err := a.repo.GetAccountByID(ctx, ID)
	if err != nil {
		return nil, err
	}

	if !accountExists {
		return nil, ErrAccountNotFound
	}

	schema, schemaFound, err := a.repo.GetSchemaByID(ctx, account.Schema)
	if err != nil {
		return nil, err
	}

	if !schemaFound {
		return nil, ErrSchemaNotFound
	}

	providers, err := a.repo.GetProvidersByIDs(ctx, schema.Providers.Values()...)
	if err != nil {
		return nil, err
	}

	airlinesCodes := hashset.New[entity.AirlineCode]()
	for _, provider := range providers {
		for _, code := range provider.Airlines.Values() {
			airlinesCodes.Put(code)
		}
	}

	return a.repo.GetAirlinesByCodes(ctx, airlinesCodes.Values()...)
}
