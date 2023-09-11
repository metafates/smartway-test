package usecase

import (
	"context"
	"errors"

	"github.com/metafates/smartway-test/internal/entity"
)

var (
	ErrUsedSchemaDeletion = errors.New("schema that is being used can not be deleted")
	ErrSchemaNotFound     = errors.New("schema not found")
)

var _ Schema = (*SchemaUseCase)(nil)

type SchemaUseCase struct {
	repo Repository
}

func NewSchemaUseCase(repository Repository) *SchemaUseCase {
	return &SchemaUseCase{repo: repository}
}

func (s *SchemaUseCase) Add(ctx context.Context, schema entity.Schema) error {
	return s.repo.StoreSchema(ctx, schema)
}

func (s *SchemaUseCase) Delete(ctx context.Context, ID string) error {
	accounts, err := s.repo.GetAccounts(ctx)
	if err != nil {
		return err
	}

	for _, account := range accounts {
		if account.SchemaID == ID {
			return ErrUsedSchemaDeletion
		}
	}

	return s.repo.DeleteSchema(ctx, ID)
}

func (s *SchemaUseCase) Update(ctx context.Context, ID string, changes entity.Schema) error {
	return s.repo.UpdateSchema(ctx, ID, changes)
}

func (s *SchemaUseCase) Find(ctx context.Context, name string) (entity.Schema, bool, error) {
	return s.repo.GetSchemaByName(ctx, name)
}
