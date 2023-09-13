package repository

import (
	"context"
	"strconv"
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

	query := p.buildBulkInsertQuery("schema_provider", []string{"schema_id", "provider_id"}, schema.Providers.Len())
	values := make([]any, schema.Providers.Len()*2)

	for _, provider := range schema.Providers.Values() {
		values = append(values, schema.ID, provider)
	}

	_, err = tx.Exec(ctx, query, values...)
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
			query := p.buildBulkInsertQuery("schema_provider", []string{"schema_id", "provider_id"}, changes.Providers.Len())
			values := make([]any, changes.Providers.Len()*2)

			for _, provider := range changes.Providers.Values() {
				values = append(values, ID, provider)
			}

			_, err := tx.Exec(ctx, query, values...)
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
select id, provider_id
from schema
join schema_provider a on schema.id = a.schema_id
where name = $1
`, name)
	if err != nil {
		return entity.Schema{}, false, err
	}

	schema := entity.Schema{
		Name:      name,
		Providers: hashset.New[entity.ProviderID](),
	}

	if rows.Next() {
		var (
			schemaID   entity.SchemaID
			providerID entity.ProviderID
		)

		if err := rows.Scan(&schemaID, &providerID); err != nil {
			return entity.Schema{}, false, err
		}

		schema.ID = schemaID
		schema.Providers.Put(providerID)
	} else {
		return entity.Schema{}, false, nil
	}

	for rows.Next() {
		var (
			schemaID   entity.SchemaID
			providerID entity.ProviderID
		)

		if err := rows.Scan(&schemaID, &providerID); err != nil {
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
	tx, err := p.Pool.Begin(ctx)
	if err != nil {
		return err
	}

	_, err = tx.Exec(ctx, `
insert into provider (id, name) values ($1, $2)
`, provider.ID, provider.Name)
	if err != nil {
		return err
	}

	if provider.Airlines == nil || provider.Airlines.IsEmpty() {
		return tx.Commit(ctx)
	}

	query := p.buildBulkInsertQuery("provider_airline", []string{"provider_id", "airline_code"}, provider.Airlines.Len())
	values := make([]any, provider.Airlines.Len()*2)

	for _, airline := range provider.Airlines.Values() {
		values = append(values, provider.ID, airline)
	}

	_, err = tx.Exec(ctx, query, values)
	if err != nil {
		return err
	}

	return tx.Commit(ctx)
}

func (p *PostgresRepository) UpdateProvider(ctx context.Context, ID entity.ProviderID, changes entity.ProviderChanges) error {
	if changes.Name == nil {
		return nil
	}

	_, err := p.Pool.Exec(ctx, `update provider set name = $1 where id = $2`, ID, *changes.Name)
	return err
}

func (p *PostgresRepository) GetProviderByID(ctx context.Context, ID entity.ProviderID) (entity.Provider, bool, error) {
	rows, err := p.Pool.Query(ctx, `
select name, airline_code
from provider
join provider_airline pa on provider.id = pa.provider_id
where id = $1
`, ID)
	if err != nil {
		return entity.Provider{}, false, err
	}

	provider := entity.Provider{
		ID:       ID,
		Airlines: hashset.New[entity.AirlineCode](),
	}

	if rows.Next() {
		var (
			name        string
			airlineCode entity.AirlineCode
		)

		if err := rows.Scan(&name, &airlineCode); err != nil {
			return entity.Provider{}, false, err
		}
	} else {
		return entity.Provider{}, false, nil
	}

	for rows.Next() {
		var (
			name        string
			airlineCode entity.AirlineCode
		)

		if err := rows.Scan(&name, &airlineCode); err != nil {
			return entity.Provider{}, false, err
		}

		provider.Airlines.Put(airlineCode)
	}

	return provider, true, nil
}

func (p *PostgresRepository) GetProvidersByIDs(ctx context.Context, IDs ...entity.ProviderID) ([]entity.Provider, error) {
	//TODO implement me
	panic("implement me")
}

func (p *PostgresRepository) DeleteProvider(ctx context.Context, ID entity.ProviderID) error {
	_, err := p.Pool.Exec(ctx, `delete from provider where id = $1`, ID)
	return err
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
	_, err := p.Pool.Exec(ctx, `delete from airline where code = $1`, code)
	return err
}

func NewPostgresRepository(pg *postgres.Postgres) *PostgresRepository {
	return &PostgresRepository{pg}
}
