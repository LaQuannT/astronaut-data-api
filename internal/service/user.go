package usecase

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/LaQuannT/astronaut-data-api/internal/model"
	"github.com/LaQuannT/astronaut-data-api/internal/transport/middleware"
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
	"role":     validation.Role,
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
	return u, nil
}

func (uc *userUsercase) List(ctx context.Context, limit, offset int) ([]*model.User, error) {
	requestUser, ok := ctx.Value(middleware.RequestUser).(*model.User)
	if !ok {
		return nil, errors.New("invalid request-User")
	}

	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	if requestUser.Role != model.AdminUser {
		return nil, errors.New("user is not authorised")
	}

	users, err := uc.store.List(ctx, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("error listing users: %w", err)
	}

	return users, nil
}

func (uc *userUsercase) Get(ctx context.Context, id int) (*model.User, error) {
	requestUser, ok := ctx.Value(middleware.RequestUser).(*model.User)
	if !ok {
		return nil, errors.New("invalid request-User")
	}

	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	if requestUser.ID == id {
		return requestUser, nil
	}

	if requestUser.Role != model.AdminUser {
		return nil, errors.New("user is not authorised")
	}

	u, err := uc.store.Get(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("error fetching user: %w", err)
	}

	return u, nil
}

func (uc *userUsercase) Update(ctx context.Context, u *model.User) (*model.User, []error) {
	errs := make([]error, 0)

	requestUser, ok := ctx.Value(middleware.RequestUser).(*model.User)
	if !ok {
		err := errors.New("invalid request-User")
		errs = append(errs, err)
		return nil, errs
	}

	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	var originalUser *model.User

	if requestUser.ID == u.ID {
		originalUser = requestUser
	} else if requestUser.ID != u.ID && requestUser.Role == model.AdminUser {
		ou, err := uc.store.Get(ctx, u.ID)
		if err != nil {
			err = fmt.Errorf("error fetching original user data: %w", err)
			errs = append(errs, err)
			return nil, errs

		}
		originalUser = ou

	} else {
		err := errors.New("user is not authorised")
		errs = append(errs, err)
		return nil, errs
	}

	v := validation.New(userValidatorRules)

	checks := map[string]validation.Check{
		"firstName": {Value: u.FirstName, RuleKey: []string{"require", "length"}},
		"surename":  {Value: u.Surename, RuleKey: []string{"require", "length"}},
		"email":     {Value: u.Email, RuleKey: []string{"require", "email", "length"}},
		"role":      {Value: u.Role, RuleKey: []string{"role"}},
	}

	if errs := v.Validate(checks); errs != nil {
		return nil, errs
	}

	u = compareUserData(originalUser, u)

	u.UpdatedAt = time.Now().UTC()

	if err := uc.store.Update(ctx, u); err != nil {
		err = fmt.Errorf("error updating user data: %w", err)
		errs = append(errs, err)
		return nil, errs
	}

	return u, nil
}

func (uc *userUsercase) Delete(ctx context.Context, id int) error {
	requestUser, ok := ctx.Value(middleware.RequestUser).(*model.User)
	if !ok {
		return errors.New("invalid request-User")
	}

	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	if requestUser.Role != model.AdminUser && requestUser.ID != id {
		return errors.New("user is not authorised")
	}

	if err := uc.store.Delete(ctx, id); err != nil {
		return fmt.Errorf("error deleting user: %w", err)
	}

	return nil
}

func (uc *userUsercase) ResetPassword(ctx context.Context, u *model.User) []error {
	errs := make([]error, 0)

	requestUser, ok := ctx.Value(middleware.RequestUser).(*model.User)
	if !ok {
		err := errors.New("invalid request-User")
		errs = append(errs, err)
		return errs
	}

	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	if requestUser.ID != u.ID {
		err := errors.New("user is not authorised")
		errs = append(errs, err)
		return errs

	}

	v := validation.New(userValidatorRules)

	checks := map[string]validation.Check{
		"password": {Value: u.Password, RuleKey: []string{"require", "password"}},
	}

	if errs := v.Validate(checks); errs != nil {
		return errs
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(u.Password), 12)
	if err != nil {
		err := fmt.Errorf("error generating password hash: %w", err)
		errs = append(errs, err)
		return errs

	}

	u.Password = string(hash)
	u.UpdatedAt = time.Now().UTC()

	if err := uc.store.UpdatePassword(ctx, u); err != nil {
		err = fmt.Errorf("error resetting user password: %w", err)
		errs = append(errs, err)
		return errs

	}

	return nil
}

func (uc *userUsercase) GenerateNewAPIKey(ctx context.Context, id int) (*model.User, error) {
	requestUser, ok := ctx.Value(middleware.RequestUser).(*model.User)
	if !ok {
		return nil, errors.New("invalid request-User")
	}

	if requestUser.ID != id {
		return nil, errors.New("user not authorised")
	}

	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	requestUser.ApiKey = uuid.New().String()
	requestUser.UpdatedAt = time.Now().UTC()

	if err := uc.store.UpdateAPIKey(ctx, requestUser); err != nil {
		return nil, fmt.Errorf("error generating new APIKey: %w", err)
	}

	return requestUser, nil
}

func (uc *userUsercase) SearchAPIKey(ctx context.Context, key string) (*model.User, error) {
	u, err := uc.store.SearchApiKey(ctx, key)
	if err != nil {
		return nil, fmt.Errorf("error searching for user by APIKey: %w", err)
	}

	return u, nil
}

func compareUserData(old, new *model.User) *model.User {
	if new.FirstName != old.FirstName {
		old.FirstName = new.FirstName
	}

	if new.Surename != old.Surename {
		old.Surename = new.Surename
	}

	if new.Email != old.Email {
		old.Email = new.Email
	}

	if new.Role != old.Role {
		old.Role = new.Role
	}

	return old
}
