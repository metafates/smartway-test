package usecase

import (
	"context"

	"github.com/go-playground/validator/v10"
	"github.com/metafates/smartway-test/internal/entity"
)

var _ Provider = (*ProviderUseCase)(nil)

type ProviderUseCase struct {
	repo     Repository
	validate *validator.Validate
}

func NewProviderUseCase(repository Repository) *ProviderUseCase {
	return &ProviderUseCase{
		repo:     repository,
		validate: validator.New(validator.WithRequiredStructEnabled()),
	}
}

func (p *ProviderUseCase) Add(ctx context.Context, provider entity.Provider) error {
	if err := p.validate.Struct(provider); err != nil {
		return err
	}

	return p.repo.StoreProvider(ctx, provider)
}

func (p *ProviderUseCase) Delete(ctx context.Context, ID entity.ProviderID) error {
	return p.repo.DeleteProvider(ctx, ID)
}

func (p *ProviderUseCase) GetAirlines(ctx context.Context, ID entity.ProviderID) ([]entity.Airline, error) {
	return p.repo.GetProviderAirlines(ctx, ID)
}
