package usecase

import (
	"context"
	"errors"

	"github.com/metafates/smartway-test/internal/entity"
	"github.com/metafates/smartway-test/internal/pkg/hashset"
)

var (
	ErrAirlineNotFound  = errors.New("airline not found")
	ErrEmptyAirlineCode = errors.New("airline code is empty")
	ErrEmptyAirlineName = errors.New("airline name is empty")
)

var _ Airline = (*AirlineUseCase)(nil)

type AirlineUseCase struct {
	repo Repository
}

func NewAirlineUseCase(repository Repository) *AirlineUseCase {
	return &AirlineUseCase{repo: repository}
}

func (a *AirlineUseCase) Delete(ctx context.Context, code entity.AirlineCode) error {
	return a.repo.DeleteAirline(ctx, code)
}

func (a *AirlineUseCase) Add(ctx context.Context, airline entity.Airline) error {
	if airline.Code == "" {
		return ErrEmptyAirlineCode
	}

	if airline.Name == "" {
		return ErrEmptyAirlineName
	}

	return a.repo.StoreAirline(ctx, airline)
}

func (a *AirlineUseCase) SetProviders(ctx context.Context, code entity.AirlineCode, providers []entity.ProviderID) error {
	set := hashset.New[entity.ProviderID]()
	set.PutAll(providers)

	return a.repo.UpdateAirline(ctx, code, entity.AirlineChanges{
		Providers: set,
	})
}
