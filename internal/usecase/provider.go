package usecase

import (
	"context"
	"errors"

	"github.com/metafates/smartway-test/internal/entity"
)

var (
	ErrProviderNotFound = errors.New("provider not found")
)

var _ Provider = (*ProviderUseCase)(nil)

type ProviderUseCase struct {
	repo Repository
}

func NewProviderUseCase(repository Repository) *ProviderUseCase {
	return &ProviderUseCase{repo: repository}
}

func (p ProviderUseCase) Add(ctx context.Context, provider entity.Provider) error {
	return p.repo.StoreProvider(ctx, provider)
}

func (p ProviderUseCase) Delete(ctx context.Context, ID entity.ProviderID) error {
	return p.repo.DeleteProvider(ctx, ID)
}

func (p ProviderUseCase) GetAirlines(ctx context.Context, ID entity.ProviderID) ([]entity.Airline, error) {
	provider, providerFound, err := p.repo.GetProviderByID(ctx, ID)
	if err != nil {
		return nil, err
	}

	if !providerFound {
		return nil, ErrProviderNotFound
	}

	return p.repo.GetAirlinesByCodes(ctx, provider.Airlines.Values()...)
}
