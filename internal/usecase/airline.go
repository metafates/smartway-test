package usecase

import (
	"context"
	"errors"

	"github.com/metafates/smartway-test/internal/entity"
	"golang.org/x/exp/maps"
)

var (
	ErrAirlineNotFound = errors.New("airline not found")
)

var _ Airline = (*AirlineUseCase)(nil)

type AirlineUseCase struct {
	repo Repository
}

func NewAirlineUseCase(repository Repository) *AirlineUseCase {
	return &AirlineUseCase{repo: repository}
}

func (a *AirlineUseCase) Add(ctx context.Context, airline entity.Airline) error {
	return a.repo.StoreAirline(ctx, airline)
}

func (a *AirlineUseCase) SetProviders(ctx context.Context, airlineCode string, providersIDs []int) error {
	airline, airlineFound, err := a.repo.GetAirlineByCode(ctx, airlineCode)
	if err != nil {
		return err
	}

	if !airlineFound {
		return ErrAirlineNotFound
	}

	providers, err := a.repo.GetProvidersByIDs(ctx, providersIDs...)
	if err != nil {
		return err
	}

	// TODO!: use transactions
	//
	// maybe this?
	// https://www.conf42.com/Golang_2023_Ilia_Sergunin_transaction_management_repository_pattern
	maps.Clear(airline.ProvidersIDs)
	for _, provider := range providers {
		airline.ProvidersIDs[provider.ID] = struct{}{}
		provider.AirlinesCodes[airline.Code] = struct{}{}

		err := a.repo.UpdateProvider(ctx, provider.ID, entity.Provider{
			AirlinesCodes: provider.AirlinesCodes,
		})

		if err != nil {
			return err
		}
	}

	return nil
}
