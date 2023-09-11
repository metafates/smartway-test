package usecase

import (
	"context"

	"github.com/metafates/smartway-test/internal/entity"
	"github.com/samber/lo"
)

var _ Provider = (*ProviderUseCase)(nil)

type ProviderUseCase struct {
	repo Repository
}

func (p ProviderUseCase) Add(ctx context.Context, provider entity.Provider) error {
	return p.repo.StoreProvider(ctx, provider, false)
}

func (p ProviderUseCase) Delete(ctx context.Context, ID string) error {
	return p.Delete(ctx, ID)
}

func (p ProviderUseCase) GetAirlines(ctx context.Context, ID string) ([]entity.Airline, error) {
	provider, providerFound, err := p.repo.GetProviderByID(ctx, ID)
	if err != nil {
		return nil, err
	}

	if !providerFound {
		return nil, ErrProviderNotFound
	}

	return p.repo.GetAirlinesByCodes(ctx, lo.Keys(provider.Airlines)...)
}
