package repository

import (
	"context"
	"errors"

	"github.com/metafates/smartway-test/internal/entity"
	"github.com/metafates/smartway-test/internal/pkg/hashset"
	"github.com/metafates/smartway-test/internal/usecase"
	"github.com/samber/lo"
)

var _ usecase.Repository = (*MockRepository)(nil)

type MockRepository struct {
	accounts  map[entity.AccountID]entity.Account
	schemas   map[entity.SchemaID]entity.Schema
	providers map[entity.ProviderID]entity.Provider
	airlines  map[entity.AirlineCode]entity.Airline
}

func NewMockRepository() *MockRepository {
	return &MockRepository{
		accounts:  make(map[entity.AccountID]entity.Account),
		schemas:   make(map[entity.SchemaID]entity.Schema),
		providers: make(map[entity.ProviderID]entity.Provider),
		airlines:  make(map[entity.AirlineCode]entity.Airline),
	}
}

func (m *MockRepository) StoreAccount(ctx context.Context, account entity.Account) error {
	if _, ok := m.accounts[account.ID]; ok {
		return errors.New("account exists")
	}

	m.accounts[account.ID] = account
	return nil
}

func (m *MockRepository) GetAccountByID(ctx context.Context, ID entity.AccountID) (entity.Account, bool, error) {
	account, ok := m.accounts[ID]
	return account, ok, nil
}

func (m *MockRepository) GetAccounts(ctx context.Context) ([]entity.Account, error) {
	return lo.Values(m.accounts), nil
}

func (m *MockRepository) DeleteAccount(ctx context.Context, ID entity.AccountID) error {
	if _, ok := m.accounts[ID]; !ok {
		return errors.New("account does not exist")
	}

	delete(m.accounts, ID)
	return nil
}

func (m *MockRepository) StoreSchema(ctx context.Context, schema entity.Schema) error {
	if _, ok := m.schemas[schema.ID]; ok {
		return errors.New("schema exists")
	}

	if schema.Providers == nil {
		schema.Providers = hashset.New[entity.ProviderID]()
	}

	m.schemas[schema.ID] = schema
	return nil
}

func (m *MockRepository) UpdateSchema(ctx context.Context, ID entity.SchemaID, changes entity.Schema) error {
	schema, ok := m.schemas[ID]
	if !ok {
		return errors.New("schema not found")
	}

	if changes.Name != "" {
		schema.Name = changes.Name
	}

	if changes.ID != 0 {
		schema.ID = changes.ID
	}

	if changes.Providers != nil {
		schema.Providers = changes.Providers
	}

	delete(m.schemas, ID)
	m.schemas[schema.ID] = schema
	return nil
}

func (m *MockRepository) GetSchemaByID(ctx context.Context, ID entity.SchemaID) (entity.Schema, bool, error) {
	schema, ok := m.schemas[ID]
	return schema, ok, nil
}

func (m *MockRepository) GetSchemaByName(ctx context.Context, name string) (entity.Schema, bool, error) {
	schema, ok := lo.Find(lo.Values(m.schemas), func(schema entity.Schema) bool {
		return schema.Name == name
	})

	return schema, ok, nil
}

func (m *MockRepository) DeleteSchema(ctx context.Context, ID entity.SchemaID) error {
	if _, ok := m.schemas[ID]; !ok {
		return errors.New("schema does not exist")
	}

	delete(m.schemas, ID)
	return nil
}

func (m *MockRepository) StoreProvider(ctx context.Context, provider entity.Provider) error {
	if _, ok := m.providers[provider.ID]; ok {
		return errors.New("provider exists")
	}

	if provider.Airlines == nil {
		provider.Airlines = hashset.New[entity.AirlineCode]()
	}

	m.providers[provider.ID] = provider
	return nil
}

func (m *MockRepository) UpdateProvider(ctx context.Context, ID entity.ProviderID, changes entity.Provider) error {
	provider, ok := m.providers[ID]
	if !ok {
		return errors.New("provider not found")
	}

	if changes.Name != "" {
		provider.Name = changes.Name
	}

	if changes.ID != "" {
		provider.ID = changes.ID
	}

	if changes.Airlines != nil {
		provider.Airlines = changes.Airlines
	}

	delete(m.providers, ID)
	m.providers[provider.ID] = provider
	return nil
}

func (m *MockRepository) GetProviderByID(ctx context.Context, ID entity.ProviderID) (entity.Provider, bool, error) {
	provider, ok := m.providers[ID]
	return provider, ok, nil
}

func (m *MockRepository) GetProvidersByIDs(ctx context.Context, IDs ...entity.ProviderID) ([]entity.Provider, error) {
	providers := make([]entity.Provider, len(IDs))

	for i, ID := range IDs {
		provider, ok := m.providers[ID]
		if !ok {
			return nil, errors.New("provider not found")
		}

		providers[i] = provider
	}

	return providers, nil
}

func (m *MockRepository) DeleteProvider(ctx context.Context, ID entity.ProviderID) error {
	if _, ok := m.providers[ID]; !ok {
		return errors.New("provider does not exist")
	}

	delete(m.providers, ID)
	return nil
}

func (m *MockRepository) StoreAirline(ctx context.Context, airline entity.Airline) error {
	if _, ok := m.airlines[airline.Code]; ok {
		return errors.New("airline exists")
	}

	if airline.Providers == nil {
		airline.Providers = hashset.New[entity.ProviderID]()
	}

	m.airlines[airline.Code] = airline
	return nil
}

func (m *MockRepository) GetAirlineByCode(ctx context.Context, code entity.AirlineCode) (entity.Airline, bool, error) {
	airline, ok := m.airlines[code]
	if !ok {
		return entity.Airline{}, false, nil
	}

	return airline, true, nil
}

func (m *MockRepository) GetAirlinesByCodes(ctx context.Context, codes ...entity.AirlineCode) ([]entity.Airline, error) {
	airlines := make([]entity.Airline, len(codes))

	for i, code := range codes {
		airline, ok := m.airlines[code]
		if !ok {
			return nil, errors.New("airline not found")
		}

		airlines[i] = airline
	}

	return airlines, nil
}

func (m *MockRepository) DeleteAirline(ctx context.Context, code entity.AirlineCode) error {
	if _, ok := m.airlines[code]; !ok {
		return errors.New("airline does not exist")
	}

	delete(m.airlines, code)
	return nil
}
