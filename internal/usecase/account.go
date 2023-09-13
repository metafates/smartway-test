package usecase

import (
	"context"
	"errors"

	"github.com/go-playground/validator/v10"
	"github.com/metafates/smartway-test/internal/entity"
	"github.com/metafates/smartway-test/internal/pkg/hashset"
)

var _ Account = (*AccountUseCase)(nil)

type AccountUseCase struct {
	repo     Repository
	validate *validator.Validate
}

func NewAccountUseCase(repository Repository) *AccountUseCase {
	return &AccountUseCase{
		repo:     repository,
		validate: validator.New(validator.WithRequiredStructEnabled()),
	}
}

func (a *AccountUseCase) Add(ctx context.Context, account entity.Account) error {
	if err := a.validate.Struct(account); err != nil {
		return err
	}

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

	var accountAirlines []entity.Airline

	codes := hashset.New[entity.AirlineCode]()
	for _, provider := range providers {
		// TODO: make it a bulk operation in the repo?
		providerAirlines, err := a.repo.GetProviderAirlines(ctx, provider.ID)
		if err != nil {
			return nil, err
		}

		for _, airline := range providerAirlines {
			if codes.Has(airline.Code) {
				continue
			}

			codes.Put(airline.Code)
			accountAirlines = append(accountAirlines, airline)
		}
	}

	return accountAirlines, nil
}
