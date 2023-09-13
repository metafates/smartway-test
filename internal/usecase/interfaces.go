package usecase

import (
	"context"

	"github.com/metafates/smartway-test/internal/entity"
)

type UseCases struct {
	Account  Account
	Schema   Schema
	Provider Provider
	Airline  Airline
}

type (
	Account interface {
		Add(ctx context.Context, account entity.Account) error
		Delete(ctx context.Context, ID entity.AccountID) error
		SetSchema(ctx context.Context, accountID entity.AccountID, schemaID entity.SchemaID) error
		GetAirlines(ctx context.Context, ID entity.AccountID) ([]entity.Airline, error)
	}

	Schema interface {
		Add(ctx context.Context, schema entity.Schema) error
		Delete(ctx context.Context, ID entity.SchemaID) error
		Update(ctx context.Context, ID entity.SchemaID, changes entity.SchemaChanges) error
		Find(ctx context.Context, name string) (entity.Schema, bool, error)
		GetProviders(ctx context.Context, ID entity.SchemaID) ([]entity.Provider, error)
	}

	Provider interface {
		Add(ctx context.Context, provider entity.Provider) error
		Delete(ctx context.Context, ID entity.ProviderID) error
		GetAirlines(ctx context.Context, ID entity.ProviderID) ([]entity.Airline, error)
	}

	Airline interface {
		Add(ctx context.Context, airline entity.Airline) error
		Delete(ctx context.Context, code entity.AirlineCode) error
		SetProviders(ctx context.Context, airlineCode entity.AirlineCode, providersIDs []entity.ProviderID) error
	}

	Repository interface {
		StoreAccount(ctx context.Context, account entity.Account) error
		DeleteAccount(ctx context.Context, ID entity.AccountID) error
		UpdateAccount(ctx context.Context, ID entity.AccountID, changes entity.AccountChanges) error
		GetAccountByID(ctx context.Context, ID entity.AccountID) (entity.Account, bool, error)
		GetAccountSchema(ctx context.Context, ID entity.AccountID) (entity.Schema, bool, error)

		StoreSchema(ctx context.Context, schema entity.Schema) error
		DeleteSchema(ctx context.Context, ID entity.SchemaID) error
		UpdateSchema(ctx context.Context, ID entity.SchemaID, changes entity.SchemaChanges) error
		GetSchemaAccounts(ctx context.Context, ID entity.SchemaID) ([]entity.Account, error)
		GetSchemaProviders(ctx context.Context, ID entity.SchemaID) ([]entity.Provider, error)
		GetSchemaByName(ctx context.Context, name string) (entity.Schema, bool, error)

		StoreProvider(ctx context.Context, provider entity.Provider) error
		DeleteProvider(ctx context.Context, ID entity.ProviderID) error
		GetProviderAirlines(ctx context.Context, ID entity.ProviderID) ([]entity.Airline, error)

		StoreAirline(ctx context.Context, airline entity.Airline) error
		DeleteAirline(ctx context.Context, code entity.AirlineCode) error
		UpdateAirline(ctx context.Context, code entity.AirlineCode, changes entity.AirlineChanges) error
	}
)
