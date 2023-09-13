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
	schema, hasSchema, err := a.repo.GetAccountSchema(ctx, ID)
	if err != nil {
		return nil, err
	}

	if !hasSchema {
		return nil, errors.New("account does not have a schema")
	}

	providers, err := a.repo.GetSchemaProviders(ctx, schema.ID)
	if err != nil {
		return nil, err
	}

	var airlines []entity.Airline

	codes := hashset.New[entity.AirlineCode]()
	for _, provider := range providers {
		// TODO: make it a bulk operation in the repo?
		airlines, err := a.repo.GetProviderAirlines(ctx, provider.ID)
		if err != nil {
			return nil, err
		}

		for _, airline := range airlines {
			if codes.Has(airline.Code) {
				continue
			}

			codes.Put(airline.Code)
			airlines = append(airlines, airline)
		}
	}

	return airlines, nil
}
