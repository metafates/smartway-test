package repository

import (
	"context"
	"fmt"
	"strings"

	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/metafates/smartway-test/internal/entity"
	"github.com/metafates/smartway-test/internal/pkg/hashset"
	"github.com/metafates/smartway-test/internal/usecase"
	"github.com/metafates/smartway-test/pkg/postgres"
)

var _ usecase.Repository = (*PostgresRepository)(nil)

type PostgresRepository struct {
	*postgres.Postgres
}

func (p *PostgresRepository) StoreAccount(ctx context.Context, account entity.Account) error {
	return p.Pool.AcquireFunc(ctx, func(conn *pgxpool.Conn) error {
		_, err := conn.Exec(ctx, `
insert into account values ($1);
`, account.ID)

		return err
	})
}

func (p *PostgresRepository) GetAccountByID(ctx context.Context, ID entity.AccountID) (entity.Account, bool, error) {
	row := p.Pool.QueryRow(ctx, `
select id, schema_id
from account
join account_schema a on account.id = a.account_id
where id = $1;
`, ID)

	var (
		accountID entity.AccountID
		schemaID  entity.SchemaID
	)
	if err := row.Scan(&accountID, &schemaID); err != nil {
		return entity.Account{}, false, err
	}

	return entity.Account{
		ID:     accountID,
		Schema: schemaID,
	}, true, nil
}

func (p *PostgresRepository) GetAccounts(ctx context.Context) ([]entity.Account, error) {
	rows, err := p.Pool.Query(ctx, `
select id, schema_id
from account
join account_schema a on account.id = a.account_id;
`)
	if err != nil {
		return nil, err
	}

	var accounts []entity.Account
	for rows.Next() {
		var (
			accountID entity.AccountID
			schemaID  entity.SchemaID
		)

		if err := rows.Scan(&accountID, &schemaID); err != nil {
			return nil, err
		}

		account := entity.Account{
			ID:     accountID,
			Schema: schemaID,
		}

		accounts = append(accounts, account)
	}

	return accounts, nil
}

func (p *PostgresRepository) UpdateAccount(ctx context.Context, ID entity.AccountID, changes entity.AccountChanges) error {
	if changes.Schema == nil {
		return nil
	}

	_, err := p.Pool.Exec(ctx, `
update account_schema set schema_id = $1
where account_id = $2;
`, changes.Schema, ID)

	return err
}

func (p *PostgresRepository) DeleteAccount(ctx context.Context, ID entity.AccountID) error {
	_, err := p.Pool.Exec(ctx, `
delete from account
where id = $1;
`, ID)

	return err
}

func (p *PostgresRepository) StoreSchema(ctx context.Context, schema entity.Schema) error {
	tx, err := p.Pool.Begin(ctx)
	if err != nil {
		return err
	}

	_, err = tx.Exec(ctx, `
insert into schema (name, id) 
values ($1, $2);
`, schema.Name, schema.ID)
	if err != nil {
		return err
	}

	if schema.Providers == nil || schema.Providers.IsEmpty() {
		return tx.Commit(ctx)
	}

	var (
		query          strings.Builder
		values         []any
		placeholderIdx = 1
	)

	query.WriteString(`insert into schema_provider (schema_id, provider_id) values `)

	for _, provider := range schema.Providers.Values() {
		query.WriteString(fmt.Sprintf("($%d, $%d),", placeholderIdx, placeholderIdx+1))
		placeholderIdx += 2

		values = append(values, schema.ID, provider)
	}

	_, err = tx.Exec(ctx, strings.TrimSuffix(query.String(), ","), values...)
	if err != nil {
		return err
	}

	return tx.Commit(ctx)
}

func (p *PostgresRepository) UpdateSchema(ctx context.Context, ID entity.SchemaID, changes entity.SchemaChanges) error {
	tx, err := p.Pool.Begin(ctx)
	if err != nil {
		return err
	}

	if changes.Name != nil {
		_, err := tx.Exec(ctx, `update schema set name = $1 where id = $2`, *changes.Name, ID)
		if err != nil {
			return err
		}
	}

	if changes.Providers != nil {
		_, err := tx.Exec(ctx, `delete from schema_provider where schema_id = $1;`, ID)
		if err != nil {
			return err
		}

		if !changes.Providers.IsEmpty() {
			var (
				query          strings.Builder
				values         []any
				placeholderIdx = 1
			)

			query.WriteString("insert into schema_provider (schema_id, provider_id) values ")

			for _, provider := range changes.Providers.Values() {
				query.WriteString(fmt.Sprintf("($%d, $%d),", placeholderIdx, placeholderIdx+1))
				placeholderIdx += 2

				values = append(values, ID, provider)
			}

			_, err := tx.Exec(ctx, strings.TrimSuffix(query.String(), ","), values...)
			if err != nil {
				return err
			}
		}
	}

	return tx.Commit(ctx)
}

func (p *PostgresRepository) GetSchemaByID(ctx context.Context, ID entity.SchemaID) (entity.Schema, bool, error) {
	rows, err := p.Pool.Query(ctx, `
select name, id, provider_id
from schema
join schema_provider a on schema.id = a.schema_id
where schema_id = $1
`, ID)
	if err != nil {
		return entity.Schema{}, false, err
	}

	schema := entity.Schema{
		Providers: hashset.New[entity.ProviderID](),
	}

	if rows.Next() {
		var (
			name       string
			schemaID   entity.SchemaID
			providerID entity.ProviderID
		)

		if err := rows.Scan(&name, &schemaID, &providerID); err != nil {
			return entity.Schema{}, false, err
		}

		schema.ID = schemaID
		schema.Name = name
		schema.Providers.Put(providerID)
	} else {
		return entity.Schema{}, false, nil
	}

	for rows.Next() {
		var (
			name       string
			schemaID   entity.SchemaID
			providerID entity.ProviderID
		)

		if err := rows.Scan(&name, &schemaID, &providerID); err != nil {
			return entity.Schema{}, false, err
		}

		schema.Providers.Put(providerID)
	}

	return schema, true, nil
}

func (p *PostgresRepository) GetSchemaByName(ctx context.Context, name string) (entity.Schema, bool, error) {
	rows, err := p.Pool.Query(ctx, `
select name, id, provider_id
from schema
join schema_provider a on schema.id = a.schema_id
where name = $1
`, name)
	if err != nil {
		return entity.Schema{}, false, err
	}

	schema := entity.Schema{
		Providers: hashset.New[entity.ProviderID](),
	}

	if rows.Next() {
		var (
			name       string
			schemaID   entity.SchemaID
			providerID entity.ProviderID
		)

		if err := rows.Scan(&name, &schemaID, &providerID); err != nil {
			return entity.Schema{}, false, err
		}

		schema.ID = schemaID
		schema.Name = name
		schema.Providers.Put(providerID)
	} else {
		return entity.Schema{}, false, nil
	}

	for rows.Next() {
		var (
			name       string
			schemaID   entity.SchemaID
			providerID entity.ProviderID
		)

		if err := rows.Scan(&name, &schemaID, &providerID); err != nil {
			return entity.Schema{}, false, err
		}

		schema.Providers.Put(providerID)
	}

	return schema, true, nil
}

func (p *PostgresRepository) DeleteSchema(ctx context.Context, ID entity.SchemaID) error {
	_, err := p.Pool.Exec(ctx, `delete from schema where id = $1`, ID)
	return err
}

func (p *PostgresRepository) StoreProvider(ctx context.Context, provider entity.Provider) error {
	//TODO implement me
	panic("implement me")
}

func (p *PostgresRepository) UpdateProvider(ctx context.Context, ID entity.ProviderID, changes entity.ProviderChanges) error {
	//TODO implement me
	panic("implement me")
}

func (p *PostgresRepository) GetProviderByID(ctx context.Context, ID entity.ProviderID) (entity.Provider, bool, error) {
	//TODO implement me
	panic("implement me")
}

func (p *PostgresRepository) GetProvidersByIDs(ctx context.Context, IDs ...entity.ProviderID) ([]entity.Provider, error) {
	//TODO implement me
	panic("implement me")
}

func (p *PostgresRepository) DeleteProvider(ctx context.Context, ID entity.ProviderID) error {
	//TODO implement me
	panic("implement me")
}

func (p *PostgresRepository) StoreAirline(ctx context.Context, airline entity.Airline) error {
	//TODO implement me
	panic("implement me")
}

func (p *PostgresRepository) UpdateAirline(ctx context.Context, code entity.AirlineCode, changes entity.AirlineChanges) error {
	//TODO implement me
	panic("implement me")
}

func (p *PostgresRepository) GetAirlineByCode(ctx context.Context, code entity.AirlineCode) (entity.Airline, bool, error) {
	//TODO implement me
	panic("implement me")
}

func (p *PostgresRepository) GetAirlinesByCodes(ctx context.Context, codes ...entity.AirlineCode) ([]entity.Airline, error) {
	//TODO implement me
	panic("implement me")
}

func (p *PostgresRepository) DeleteAirline(ctx context.Context, code entity.AirlineCode) error {
	//TODO implement me
	panic("implement me")
}

func NewPostgresRepository(pg *postgres.Postgres) *PostgresRepository {
	return &PostgresRepository{pg}
}
