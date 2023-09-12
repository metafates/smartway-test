package repository

import (
	"context"

	"github.com/metafates/smartway-test/internal/entity"
	"github.com/metafates/smartway-test/internal/usecase"
	"github.com/metafates/smartway-test/pkg/postgres"
)

var _ usecase.Repository = (*PostgresRepository)(nil)

type PostgresRepository struct {
	*postgres.Postgres
}

func (p PostgresRepository) StoreAccount(ctx context.Context, account entity.Account) error {
	//TODO implement me
	panic("implement me")
}

func (p PostgresRepository) GetAccountByID(ctx context.Context, ID int) (entity.Account, bool, error) {
	//TODO implement me
	panic("implement me")
}

func (p PostgresRepository) GetAccounts(ctx context.Context) ([]entity.Account, error) {
	//TODO implement me
	panic("implement me")
}

func (p PostgresRepository) DeleteAccount(ctx context.Context, ID int) error {
	//TODO implement me
	panic("implement me")
}

func (p PostgresRepository) StoreSchema(ctx context.Context, schema entity.Schema) error {
	//TODO implement me
	panic("implement me")
}

func (p PostgresRepository) UpdateSchema(ctx context.Context, ID int, changes entity.Schema) error {
	//TODO implement me
	panic("implement me")
}

func (p PostgresRepository) GetSchemaByID(ctx context.Context, ID int) (entity.Schema, bool, error) {
	//TODO implement me
	panic("implement me")
}

func (p PostgresRepository) GetSchemaByName(ctx context.Context, name string) (entity.Schema, bool, error) {
	//TODO implement me
	panic("implement me")
}

func (p PostgresRepository) DeleteSchema(ctx context.Context, ID int) error {
	//TODO implement me
	panic("implement me")
}

func (p PostgresRepository) StoreProvider(ctx context.Context, provider entity.Provider) error {
	//TODO implement me
	panic("implement me")
}

func (p PostgresRepository) UpdateProvider(ctx context.Context, ID int, changes entity.Provider) error {
	//TODO implement me
	panic("implement me")
}

func (p PostgresRepository) GetProviderByID(ctx context.Context, ID int) (entity.Provider, bool, error) {
	//TODO implement me
	panic("implement me")
}

func (p PostgresRepository) GetProvidersByIDs(ctx context.Context, IDs ...int) ([]entity.Provider, error) {
	//TODO implement me
	panic("implement me")
}

func (p PostgresRepository) DeleteProvider(ctx context.Context, ID int) error {
	//TODO implement me
	panic("implement me")
}

func (p PostgresRepository) StoreAirline(ctx context.Context, airline entity.Airline) error {
	//TODO implement me
	panic("implement me")
}

func (p PostgresRepository) GetAirlineByCode(ctx context.Context, code string) (entity.Airline, bool, error) {
	//TODO implement me
	panic("implement me")
}

func (p PostgresRepository) GetAirlinesByCodes(ctx context.Context, codes ...string) ([]entity.Airline, error) {
	//TODO implement me
	panic("implement me")
}
