package usecase

import (
	"context"
	"fmt"

	"github.com/LaQuannT/astronaut-data-api/internal/model"
)

type astronautUsecase struct {
	astronautStore model.AstronautStore
	userStore      model.UserStore
}

func NewAstronautUsecase(as model.AstronautStore, us model.UserStore) *astronautUsecase {
	return &astronautUsecase{
		astronautStore: as,
		userStore:      us,
	}
}

func (uc *astronautUsecase) Create(ctx context.Context, a *model.Astronaut) (*model.Astronaut, error) {
	// TODO - validate apikey and admin permission, validate astronaut data

	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	id, err := uc.astronautStore.Create(ctx, a)
	if err != nil {
		return nil, fmt.Errorf("error creating new astronaut: %w", err)
	}

	a.ID = id
	return a, nil
}

func (uc *astronautUsecase) List(ctx context.Context, limit, offset int) ([]*model.Astronaut, error) {
	// TODO- validate apikey base user permission

	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	astronauts, err := uc.astronautStore.List(ctx, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("error listing astronauts: %w", err)
	}

	return astronauts, nil
}

func (uc *astronautUsecase) Get(ctx context.Context, id int) (*model.Astronaut, error) {
	// TODO - validate apikey base user permission

	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	a, err := uc.astronautStore.Get(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("error fetching astronaut data: %w", err)
	}

	return a, nil
}

func (uc *astronautUsecase) Update(ctx context.Context, a *model.Astronaut) (*model.Astronaut, error) {
	// TODO - validate apikey and admin permission

	_, err := uc.astronautStore.Get(ctx, a.ID)
	if err != nil {
		return nil, fmt.Errorf("error fetching original astronaut data: %w", err)
	}

	// TODO - validate and compare original and new astronaut data

	if err := uc.astronautStore.Update(ctx, a); err != nil {
		return nil, fmt.Errorf("error updating astronaut data: %w", err)
	}

	return a, nil
}

func (uc *astronautUsecase) Delete(ctx context.Context, id int) error {
	// TODO - validate apikey and admin permission

	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	if err := uc.astronautStore.Delete(ctx, id); err != nil {
		return fmt.Errorf("error deleting astronaut: %w", err)
	}

	return nil
}
