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

	Scheme interface {
		Add(ctx context.Context, scheme entity.Schema) error
		Delete(ctx context.Context, ID string) error
		Change(ctx context.Context, scheme, changes entity.Schema) (entity.Schema, error)
		Find(ctx context.Context, name string) (entity.Schema, bool, error)
	}

	Provider interface {
		Add(ctx context.Context, supplier entity.Provider) error
		Delete(ctx context.Context, ID string) error
		GetAirlines(ctx context.Context, ID string) ([]entity.Airline, error)
	}

	Airline interface {
		Add(ctx context.Context, airline entity.Airline) error
		SetProviders(ctx context.Context, airlineID string, providersCodes []string) error
	}

	Repository interface {
		StoreAccount(ctx context.Context, account entity.Account, overwrite bool) error
		GetAccount(ctx context.Context, ID string) (entity.Account, bool, error)
		DeleteAccount(ctx context.Context, ID string) error

		StoreSchema(ctx context.Context, scheme entity.Schema, overwrite bool) error
		GetSchema(ctx context.Context, schemeID string) (entity.Schema, bool, error)

		StoreProvider(ctx context.Context, provider entity.Provider, overwrite bool) error
		GetProvider(ctx context.Context, ID string) (entity.Provider, bool, error)
		GetProviders(ctx context.Context, IDs ...string) ([]entity.Provider, error)

		StoreAirline(ctx context.Context, airline entity.Airline, overwrite bool) error
		GetAirlines(ctx context.Context, IDs ...string) ([]entity.Airline, error)
	}
)
