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
		GetAccountByID(ctx context.Context, ID entity.AccountID) (entity.Account, bool, error)
		GetAccounts(ctx context.Context) ([]entity.Account, error)
		UpdateAccount(ctx context.Context, ID entity.AccountID, changes entity.AccountChanges) error
		DeleteAccount(ctx context.Context, ID entity.AccountID) error

		StoreSchema(ctx context.Context, schema entity.Schema) error
		UpdateSchema(ctx context.Context, ID entity.SchemaID, changes entity.SchemaChanges) error
		GetSchemaByID(ctx context.Context, ID entity.SchemaID) (entity.Schema, bool, error)
		GetSchemaByName(ctx context.Context, name string) (entity.Schema, bool, error)
		DeleteSchema(ctx context.Context, ID entity.SchemaID) error

		StoreProvider(ctx context.Context, provider entity.Provider) error
		UpdateProvider(ctx context.Context, ID entity.ProviderID, changes entity.ProviderChanges) error
		GetProviderByID(ctx context.Context, ID entity.ProviderID) (entity.Provider, bool, error)
		GetProvidersByIDs(ctx context.Context, IDs ...entity.ProviderID) ([]entity.Provider, error)
		DeleteProvider(ctx context.Context, ID entity.ProviderID) error

		StoreAirline(ctx context.Context, airline entity.Airline) error
		UpdateAirline(ctx context.Context, code entity.AirlineCode, changes entity.AirlineChanges) error
		GetAirlineByCode(ctx context.Context, code entity.AirlineCode) (entity.Airline, bool, error)
		GetAirlinesByCodes(ctx context.Context, codes ...entity.AirlineCode) ([]entity.Airline, error)
		DeleteAirline(ctx context.Context, code entity.AirlineCode) error
	}
)
