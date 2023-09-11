package usecase

import (
	"context"

	"github.com/metafates/smartway-test/internal/entity"
)

var _ Schema = (*SchemaUseCase)(nil)

type SchemaUseCase struct {
	repo Repository
}

func (s *SchemaUseCase) Add(ctx context.Context, schema entity.Schema) error {
	return s.repo.StoreSchema(ctx, schema, false)
}

func (s *SchemaUseCase) Delete(ctx context.Context, ID string) error {
	return s.repo.DeleteSchema(ctx, ID)
}

func (s *SchemaUseCase) Update(ctx context.Context, ID string, changes entity.Schema) (entity.Schema, error) {
	schema, schemaFound, err := s.repo.GetSchemaByID(ctx, ID)
	if err != nil {
		return entity.Schema{}, err
	}

	if !schemaFound {
		return entity.Schema{}, ErrSchemaNotFound
	}

	updated := s.update(schema, changes)

	err = s.repo.StoreSchema(ctx, updated, true)
	if err != nil {
		return entity.Schema{}, err
	}

	return updated, nil
}

func (s *SchemaUseCase) update(schema, changes entity.Schema) entity.Schema {
	if changes.ID != "" {
		schema.ID = changes.ID
	}

	if changes.Providers != nil {
		schema.Providers = changes.Providers
	}

	if changes.Name != "" {
		schema.Name = changes.Name
	}

	return schema
}

func (s *SchemaUseCase) Find(ctx context.Context, name string) (entity.Schema, bool, error) {
	return s.repo.GetSchemaByName(ctx, name)
}

func NewSchemaUseCase(repository Repository) *SchemaUseCase {
	return &SchemaUseCase{repo: repository}
}
