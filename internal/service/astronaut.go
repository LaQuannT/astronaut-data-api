package usecase

import (
	"context"
	"errors"
	"fmt"

	"github.com/LaQuannT/astronaut-data-api/internal/model"
	"github.com/LaQuannT/astronaut-data-api/internal/transport/middleware"
	"github.com/LaQuannT/astronaut-data-api/internal/validation"
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

var astronautValidatorRules = validation.Rules{
	"require": validation.Required,
	"status":  validation.Status,
	"date":    validation.Date,
	"gender":  validation.Gender,
}

func (uc *astronautUsecase) Create(ctx context.Context, a *model.Astronaut) (*model.Astronaut, []error) {
	errs := make([]error, 0)

	requestUser, ok := ctx.Value(middleware.RequestUser).(*model.User)
	if !ok {
		err := errors.New("invalid request-user")
		errs := append(errs, err)
		return nil, errs
	}

	if requestUser.Role != model.AdminUser {
		err := errors.New("user is not authorised")
		errs := append(errs, err)
		return nil, errs
	}

	v := validation.New(astronautValidatorRules)

	checks := map[string]validation.Check{
		"name":        {Value: a.Name, RuleKey: []string{"require"}},
		"status":      {Value: a.Status, RuleKey: []string{"status"}},
		"birth date":  {Value: a.BirthDate, RuleKey: []string{"date"}},
		"birth place": {Value: a.BirthPlace, RuleKey: []string{"require"}},
		"gender":      {Value: a.Gender, RuleKey: []string{"gender"}},
	}

	if errs := v.Validate(checks); errs != nil {
		return nil, errs
	}

	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	a, err := uc.astronautStore.Create(ctx, a)
	if err != nil {
		err := fmt.Errorf("error creating new astronaut: %w", err)
		errs := append(errs, err)
		return nil, errs
	}

	return a, nil
}

func (uc *astronautUsecase) List(ctx context.Context, limit, offset int) ([]*model.Astronaut, error) {
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	astronauts, err := uc.astronautStore.List(ctx, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("error listing astronauts: %w", err)
	}

	return astronauts, nil
}

func (uc *astronautUsecase) Get(ctx context.Context, id int) (*model.Astronaut, error) {
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	a, err := uc.astronautStore.Get(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("error fetching astronaut data: %w", err)
	}

	return a, nil
}

func (uc *astronautUsecase) Update(ctx context.Context, a *model.Astronaut) (*model.Astronaut, error) {
	requestUser, ok := ctx.Value(middleware.RequestUser).(*model.User)
	if !ok {
		return nil, errors.New("invalid request-user")
	}

	if requestUser.Role != model.AdminUser {
		return nil, errors.New("user is not authorised")
	}

	original, err := uc.astronautStore.Get(ctx, a.ID)
	if err != nil {
		return nil, fmt.Errorf("error fetching original astronaut data: %w", err)
	}

	a = compareAstronautData(original, a)

	if err := uc.astronautStore.Update(ctx, a); err != nil {
		return nil, fmt.Errorf("error updating astronaut data: %w", err)
	}

	return a, nil
}

func (uc *astronautUsecase) Delete(ctx context.Context, id int) error {
	requestUser, ok := ctx.Value(middleware.RequestUser).(*model.User)
	if !ok {
		return errors.New("invalid request-user")
	}

	if requestUser.Role != model.AdminUser {
		return errors.New("user is not authorised")
	}

	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	if err := uc.astronautStore.Delete(ctx, id); err != nil {
		return fmt.Errorf("error deleting astronaut: %w", err)
	}

	return nil
}

func compareAstronautData(old, new *model.Astronaut) *model.Astronaut {
	if new.Name != "" && new.Name != old.Name {
		old.Name = new.Name
	}
	if new.Year != 0 && new.Year != old.Year {
		old.Year = new.Year
	}
	if new.Group != 0 && new.Group != old.Group {
		old.Group = new.Group
	}
	if new.Status != "" && new.Status != old.Status {
		old.Status = new.Status
	}
	if new.BirthDate != "" && new.Status != old.Status {
		old.Status = new.Status
	}
	if new.BirthPlace != "" && new.Status != old.BirthPlace {
		old.BirthPlace = new.BirthPlace
	}
	if new.Gender != "" && new.Gender != old.Gender {
		old.Gender = new.Gender
	}
	return old
}
