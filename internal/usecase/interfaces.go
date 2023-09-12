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
		Delete(ctx context.Context, ID int) error
		SetSchema(ctx context.Context, accountID int, schemaID int) error
		GetAirlines(ctx context.Context, ID int) ([]entity.Airline, error)
	}

	Schema interface {
		Add(ctx context.Context, schema entity.Schema) error
		Delete(ctx context.Context, ID int) error
		Update(ctx context.Context, ID int, changes entity.Schema) error
		Find(ctx context.Context, name string) (entity.Schema, bool, error)
	}

	Provider interface {
		Add(ctx context.Context, provider entity.Provider) error
		Delete(ctx context.Context, ID int) error
		GetAirlines(ctx context.Context, ID int) ([]entity.Airline, error)
	}

	Airline interface {
		Add(ctx context.Context, airline entity.Airline) error
		SetProviders(ctx context.Context, airlineCode string, providersIDs []int) error
	}

	Repository interface {
		StoreAccount(ctx context.Context, account entity.Account) error
		GetAccountByID(ctx context.Context, ID int) (entity.Account, bool, error)
		GetAccounts(ctx context.Context) ([]entity.Account, error)
		DeleteAccount(ctx context.Context, ID int) error

		StoreSchema(ctx context.Context, schema entity.Schema) error
		UpdateSchema(ctx context.Context, ID int, changes entity.Schema) error
		GetSchemaByID(ctx context.Context, ID int) (entity.Schema, bool, error)
		GetSchemaByName(ctx context.Context, name string) (entity.Schema, bool, error)
		DeleteSchema(ctx context.Context, ID int) error

		StoreProvider(ctx context.Context, provider entity.Provider) error
		UpdateProvider(ctx context.Context, ID int, changes entity.Provider) error
		GetProviderByID(ctx context.Context, ID int) (entity.Provider, bool, error)
		GetProvidersByIDs(ctx context.Context, IDs ...int) ([]entity.Provider, error)
		DeleteProvider(ctx context.Context, ID int) error

		StoreAirline(ctx context.Context, airline entity.Airline) error
		GetAirlineByCode(ctx context.Context, code string) (entity.Airline, bool, error)
		GetAirlinesByCodes(ctx context.Context, codes ...string) ([]entity.Airline, error)
	}
)
