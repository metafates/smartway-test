package usecase

import (
	"context"
	"errors"

	"github.com/metafates/smartway-test/internal/entity"
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
	airline.Providers.Clear()
	for _, provider := range providers {
		airline.Providers.Put(provider.ID)
		provider.Airlines.Put(airline.Code)

		err := a.repo.UpdateProvider(ctx, provider.ID, entity.Provider{
			Airlines: provider.Airlines,
		})

		if err != nil {
			return err
		}
	}

	return nil
}
