package usecase

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/LaQuannT/astronaut-data-api/internal/model"
	"github.com/LaQuannT/astronaut-data-api/internal/transport/handler"
	"github.com/LaQuannT/astronaut-data-api/internal/validation"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

var timeout = 5 * time.Second

type userUsercase struct {
	store model.UserStore
}

func NewUserUsecase(s model.UserStore) *userUsercase {
	return &userUsercase{
		store: s,
	}
}

var userValidatorRules = validation.Rules{
	"require":  validation.Required,
	"length":   validation.Length(50),
	"email":    validation.Email,
	"password": validation.Password(8),
}

func (uc *userUsercase) Create(ctx context.Context, u *model.User) (*model.User, []error) {
	errs := make([]error, 0)
	v := validation.New(userValidatorRules)

	checks := map[string]validation.Check{
		"firstName": {Value: u.FirstName, RuleKey: []string{"require", "length"}},
		"surename":  {Value: u.Surename, RuleKey: []string{"require", "length"}},
		"email":     {Value: u.Email, RuleKey: []string{"require", "email", "length"}},
		"password":  {Value: u.Password, RuleKey: []string{"require", "password"}},
	}

	if errs := v.Validate(checks); errs != nil {
		return nil, errs
	}

	if u.Role != model.AdminUser {
		u.Role = model.BaseUser
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(u.Password), 12)
	if err != nil {
		err = fmt.Errorf("error generating password hash: %w", err)
		errs = append(errs, err)
		return nil, errs
	}

	u.Password = string(hash)
	u.ApiKey = uuid.NewString()
	u.CreatedAt = time.Now().UTC()

	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	id, err := uc.store.Create(ctx, u)
	if err != nil {
		err = fmt.Errorf("error creating a new user: %w", err)
		errs = append(errs, err)
		return nil, errs
	}

	u.ID = id
	u.Password = ""
	return u, nil
}

func (uc *userUsercase) List(ctx context.Context, limit, offset int) ([]*model.User, error) {
	key := ctx.Value(handler.ApiKeyHeader)
	apikey, ok := key.(string)
	if !ok {
		return nil, errors.New("invalid API Key")
	}

	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	u, err := uc.store.SearchApiKey(ctx, apikey)
	if err != nil {
		return nil, fmt.Errorf("error searching user by api key: %w", err)
	}

	if u.Role != model.AdminUser {
		return nil, errors.New("user is not authorised")
	}

	users, err := uc.store.List(ctx, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("error listing users: %w", err)
	}

	return users, nil
}

func (uc *userUsercase) Get(ctx context.Context, id int) (*model.User, error) {
	key := ctx.Value(handler.ApiKeyHeader)
	apikey, ok := key.(string)
	if !ok {
		return nil, errors.New("invalid API Key")
	}

	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	u, err := uc.store.SearchApiKey(ctx, apikey)
	if err != nil {
		return nil, fmt.Errorf("error searching user by api key: %w", err)
	}

	if u.Role != model.AdminUser {
		return nil, errors.New("user is not authorised")
	}

	u, err = uc.store.Get(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("error fetching user: %w", err)
	}

	return u, nil
}

func (uc *userUsercase) Update(ctx context.Context, u *model.User) (*model.User, error) {
	key := ctx.Value(handler.ApiKeyHeader)
	apikey, ok := key.(string)
	if !ok {
		return nil, errors.New("invalid API Key")
	}

	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	requestUser, err := uc.store.SearchApiKey(ctx, apikey)
	if err != nil {
		return nil, fmt.Errorf("error searching user by api key: %w", err)
	}

	if requestUser.Role != model.AdminUser {
		if requestUser.ID != u.ID {
			return nil, errors.New("user is not authorised")
		}
	}

	_, err = uc.store.Get(ctx, u.ID)
	if err != nil {
		return nil, fmt.Errorf("error fetching original user data: %w", err)
	}

	// TODO - validate and compare original and new data

	u.UpdatedAt = time.Now().UTC()

	if err := uc.store.Update(ctx, u); err != nil {
		return nil, fmt.Errorf("error updating user data: %w", err)
	}

	return u, nil
}

func (uc *userUsercase) Delete(ctx context.Context, id int) error {
	key := ctx.Value(handler.ApiKeyHeader)
	apikey, ok := key.(string)
	if !ok {
		return errors.New("invalid API Key")
	}

	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	u, err := uc.store.SearchApiKey(ctx, apikey)
	if err != nil {
		return fmt.Errorf("error searching user by api key: %w", err)
	}

	if u.Role != model.AdminUser {
		if u.ID != id {
			return errors.New("user is not authorised")
		}
	}

	if err := uc.store.Delete(ctx, id); err != nil {
		return fmt.Errorf("error deleting user: %w", err)
	}

	return nil
}
