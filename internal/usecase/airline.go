package usecase

import (
	"context"

	"github.com/metafates/smartway-test/internal/entity"
)

var _ Airline = (*AirlineUseCase)(nil)

type AirlineUseCase struct {
	repo Repository
}

func (a AirlineUseCase) Add(ctx context.Context, airline entity.Airline) error {
	return a.repo.StoreAirline(ctx, airline, false)
}

func (a AirlineUseCase) SetProviders(ctx context.Context, airlineID string, providersCodes []string) error {
	//TODO implement me
	panic("implement me")
}
