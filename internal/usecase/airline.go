package usecase

import (
	"context"

	"github.com/go-playground/validator/v10"
	"github.com/metafates/smartway-test/internal/entity"
	"github.com/metafates/smartway-test/internal/pkg/hashset"
)

var _ Airline = (*AirlineUseCase)(nil)

type AirlineUseCase struct {
	repo     Repository
	validate *validator.Validate
}

func NewAirlineUseCase(repository Repository) *AirlineUseCase {
	return &AirlineUseCase{
		repo:     repository,
		validate: validator.New(validator.WithRequiredStructEnabled()),
	}
}

func (a *AirlineUseCase) Delete(ctx context.Context, code entity.AirlineCode) error {
	return a.repo.DeleteAirline(ctx, code)
}

func (a *AirlineUseCase) Add(ctx context.Context, airline entity.Airline) error {
	if err := a.validate.Struct(airline); err != nil {
		return err
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
