package repository

import (
	"context"
	"strconv"
	"strings"

	"github.com/metafates/smartway-test/internal/entity"
	"github.com/metafates/smartway-test/internal/usecase"
	"github.com/metafates/smartway-test/pkg/postgres"
)

var _ usecase.Repository = (*PostgresRepository)(nil)

type PostgresRepository struct {
	*postgres.Postgres
}

func (p *PostgresRepository) StoreAccount(ctx context.Context, account entity.Account) error {
	_, err := p.Pool.Exec(ctx, `insert into account (id) values ($1)`, account.ID)
	return err
}

func (p *PostgresRepository) DeleteAccount(ctx context.Context, ID entity.AccountID) error {
	_, err := p.Pool.Exec(ctx, `delete from account where id = $1`, ID)
	return err
}

func (p *PostgresRepository) UpdateAccount(ctx context.Context, ID entity.AccountID, changes entity.AccountChanges) error {
	if changes.Schema == nil {
		return nil
	}

	tx, err := p.Pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	_, err = tx.Exec(ctx, `delete from account_schema where account_id = $1`)
	if err != nil {
		return err
	}

	_, err = tx.Exec(ctx, `
insert into account_schema (account_id, schema_id) values ($1, $2)
`, ID, *changes.Schema)
	if err != nil {
		return err
	}

	return tx.Commit(ctx)
}

func (p *PostgresRepository) GetAccountByID(ctx context.Context, ID entity.AccountID) (entity.Account, bool, error) {
	return entity.Account{ID: ID}, false, nil
}

func (p *PostgresRepository) GetAccountSchema(ctx context.Context, ID entity.AccountID) (entity.Schema, bool, error) {
	// TODO: handle the case when account does not have a schema
	row := p.Pool.QueryRow(ctx, `
select s.id, s.name
from account_schema
join account a on a.id = account_schema.account_id
join schema s on s.id = account_schema.schema_id
where account_id = $1
`, ID)

	var (
		schemaID   entity.SchemaID
		schemaName string
	)

	if err := row.Scan(&schemaID, &schemaName); err != nil {
		return entity.Schema{}, false, err
	}

	return entity.Schema{
		Name: schemaName,
		ID:   schemaID,
	}, true, nil
}

func (p *PostgresRepository) StoreSchema(ctx context.Context, schema entity.Schema) error {
	_, err := p.Pool.Exec(ctx, `insert into schema (name, id) values ($1, $2)`, schema.Name, schema.ID)
	return err
}

func (p *PostgresRepository) DeleteSchema(ctx context.Context, ID entity.SchemaID) error {
	_, err := p.Pool.Exec(ctx, `delete from schema where id = $1`, ID)
	return err
}

func (p *PostgresRepository) UpdateSchema(ctx context.Context, ID entity.SchemaID, changes entity.SchemaChanges) error {
	tx, err := p.Pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	if changes.Name != nil {
		_, err := tx.Exec(ctx, `update schema set name = $1 where id = $2`, *changes.Name, ID)
		if err != nil {
			return err
		}
	}

	if changes.Providers != nil {
		_, err := tx.Exec(ctx, `delete from schema_provider where schema_id = $1`, ID)
		if err != nil {
			return err
		}

		query := p.buildBulkInsertQuery("schema_provider", []string{"schema_id", "provider_id"}, changes.Providers.Len())
		values := make([]any, changes.Providers.Len()*2)

		for _, provider := range changes.Providers.Values() {
			values = append(values, ID, provider)
		}

		_, err = tx.Exec(ctx, query, values)
		if err != nil {
			return err
		}
	}

	return tx.Commit(ctx)
}

func (p *PostgresRepository) GetSchemaAccounts(ctx context.Context, ID entity.SchemaID) ([]entity.Account, error) {
	rows, err := p.Pool.Query(ctx, `
select account_id
from account_schema
where schema_id = $1
`, ID)
	if err != nil {
		return nil, err
	}

	var accounts []entity.Account

	for rows.Next() {
		var (
			accountID entity.AccountID
		)

		if err := rows.Scan(&accountID); err != nil {
			return nil, err
		}

		accounts = append(accounts, entity.Account{ID: accountID})
	}

	return accounts, nil
}

func (p *PostgresRepository) GetSchemaProviders(ctx context.Context, ID entity.SchemaID) ([]entity.Provider, error) {
	rows, err := p.Pool.Query(ctx, `
select p.id, p.name
from schema_provider
join provider p on p.id = schema_provider.provider_id
where schema_id = $1
`, ID)
	if err != nil {
		return nil, err
	}

	var providers []entity.Provider

	for rows.Next() {
		var (
			providerID   entity.ProviderID
			providerName string
		)

		if err := rows.Scan(&providerID, &providerName); err != nil {
			return nil, err
		}

		providers = append(providers, entity.Provider{
			ID:   providerID,
			Name: providerName,
		})
	}

	return providers, nil
}

func (p *PostgresRepository) GetSchemaByName(ctx context.Context, name string) (entity.Schema, bool, error) {
	//TODO implement me
	panic("implement me")
}

func (p *PostgresRepository) StoreProvider(ctx context.Context, provider entity.Provider) error {
	//TODO implement me
	panic("implement me")
}

func (p *PostgresRepository) DeleteProvider(ctx context.Context, ID entity.ProviderID) error {
	//TODO implement me
	panic("implement me")
}

func (p *PostgresRepository) GetProviderAirlines(ctx context.Context, ID entity.ProviderID) ([]entity.Airline, error) {
	//TODO implement me
	panic("implement me")
}

func (p *PostgresRepository) StoreAirline(ctx context.Context, airline entity.Airline) error {
	//TODO implement me
	panic("implement me")
}

func (p *PostgresRepository) DeleteAirline(ctx context.Context, code entity.AirlineCode) error {
	//TODO implement me
	panic("implement me")
}

func (p *PostgresRepository) UpdateAirline(ctx context.Context, code entity.AirlineCode, changes entity.AirlineChanges) error {
	//TODO implement me
	panic("implement me")
}

func (p *PostgresRepository) buildBulkInsertQuery(table string, columns []string, rowCount int) string {
	var b strings.Builder
	var cnt int

	columnCount := len(columns)

	b.Grow(40000) // Need to calculate, I'm too lazy))

	b.WriteString("insert into " + table + "(" + strings.Join(columns, ", ") + ") values ")

	for i := 0; i < rowCount; i++ {
		b.WriteString("(")
		for j := 0; j < columnCount; j++ {
			cnt++
			b.WriteString("$")
			b.WriteString(strconv.Itoa(cnt))
			if j != columnCount-1 {
				b.WriteString(", ")
			}
		}
		b.WriteString(")")
		if i != rowCount-1 {
			b.WriteString(",")
		}
	}
	return b.String()
}
