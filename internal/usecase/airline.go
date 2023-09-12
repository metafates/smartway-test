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

func (a *AirlineUseCase) SetProviders(ctx context.Context, airlineCode entity.AirlineCode, providersIDs []entity.ProviderID) error {
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
	maps.Clear(airline.Providers)
	for _, provider := range providers {
		airline.Providers[provider.ID] = struct{}{}
		provider.Airlines[airline.Code] = struct{}{}

		err := a.repo.UpdateProvider(ctx, provider.ID, entity.Provider{
			Airlines: provider.Airlines,
		})

		if err != nil {
			return err
		}
	}

	return nil
}
