package usecase

import (
	"context"

	"github.com/metafates/smartway-test/internal/entity"
)

type (
	Account interface {
		Add(ctx context.Context, account entity.Account) error
		Delete(ctx context.Context, ID string) error
		SetScheme(ctx context.Context, accountID string, schemaID string) error
		GetAirlines(ctx context.Context, ID string) ([]entity.Airline, error)
	}

	Schema interface {
		Add(ctx context.Context, schema entity.Schema) error
		Delete(ctx context.Context, ID string) error
		Update(ctx context.Context, ID string, changes entity.Schema) (entity.Schema, error)
		Find(ctx context.Context, name string) (entity.Schema, bool, error)
	}

	Provider interface {
		Add(ctx context.Context, provider entity.Provider) error
		Delete(ctx context.Context, ID string) error
		GetAirlines(ctx context.Context, ID string) ([]entity.Airline, error)
	}

	Airline interface {
		Add(ctx context.Context, airline entity.Airline) error
		SetProviders(ctx context.Context, airlineID string, providersCodes []string) error
	}

	Repository interface {
		StoreAccount(ctx context.Context, account entity.Account, overwrite bool) error
		GetAccountByID(ctx context.Context, ID string) (entity.Account, bool, error)
		DeleteAccount(ctx context.Context, ID string) error

		StoreSchema(ctx context.Context, schema entity.Schema, overwrite bool) error
		GetSchemaByID(ctx context.Context, ID string) (entity.Schema, bool, error)
		GetSchemaByName(ctx context.Context, name string) (entity.Schema, bool, error)
		DeleteSchema(ctx context.Context, ID string) error

		StoreProvider(ctx context.Context, provider entity.Provider, overwrite bool) error
		GetProviderByID(ctx context.Context, ID string) (entity.Provider, bool, error)
		GetProvidersByIDs(ctx context.Context, IDs ...string) ([]entity.Provider, error)
		DeleteProvider(ctx context.Context, ID string) error

		StoreAirline(ctx context.Context, airline entity.Airline, overwrite bool) error
		GetAirlinesByCodes(ctx context.Context, codes ...string) ([]entity.Airline, error)
	}
)
