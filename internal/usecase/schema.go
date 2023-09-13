package usecase

import (
	"context"
	"errors"

	"github.com/go-playground/validator/v10"
	"github.com/metafates/smartway-test/internal/entity"
)

var (
	ErrUsedSchemaDeletion = errors.New("schema that is being used can not be deleted")
)

var _ Schema = (*SchemaUseCase)(nil)

type SchemaUseCase struct {
	repo     Repository
	validate *validator.Validate
}

func NewSchemaUseCase(repository Repository) *SchemaUseCase {
	return &SchemaUseCase{
		repo:     repository,
		validate: validator.New(validator.WithRequiredStructEnabled()),
	}
}

func (s *SchemaUseCase) Add(ctx context.Context, schema entity.Schema) error {
	if err := s.validate.Struct(schema); err != nil {
		return err
	}

	return s.repo.StoreSchema(ctx, schema)
}

func (s *SchemaUseCase) Delete(ctx context.Context, ID entity.SchemaID) error {
	accounts, err := s.repo.GetSchemaAccounts(ctx, ID)
	if err != nil {
		return err
	}

	if len(accounts) != 0 {
		return ErrUsedSchemaDeletion
	}

	return s.repo.DeleteSchema(ctx, ID)
}

func (s *SchemaUseCase) Update(ctx context.Context, ID entity.SchemaID, changes entity.SchemaChanges) error {
	return s.repo.UpdateSchema(ctx, ID, changes)
}

func (s *SchemaUseCase) Find(ctx context.Context, name string) (entity.Schema, bool, error) {
	return s.repo.GetSchemaByName(ctx, name)
}

func (s *SchemaUseCase) GetProviders(ctx context.Context, ID entity.SchemaID) ([]entity.Provider, error) {
	return s.repo.GetSchemaProviders(ctx, ID)
}
