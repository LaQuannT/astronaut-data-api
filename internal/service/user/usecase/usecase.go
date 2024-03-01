package usecase

import (
	"context"
	"fmt"
	"time"

	"github.com/LaQuannT/astronaut-data-api/internal/model"
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

func (uc *userUsercase) Create(ctx context.Context, u *model.User) (*model.User, error) {
	// TODO - validation for user data

	if u.Role != model.AdminUser {
		u.Role = model.BaseUser
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(u.Password), 12)
	if err != nil {
		return nil, fmt.Errorf("error generating password hash: %w", err)
	}

	u.Password = string(hash)
	u.ApiKey = uuid.NewString()
	u.CreatedAt = time.Now().UTC()

	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	id, err := uc.store.Create(ctx, u)
	if err != nil {
		return nil, fmt.Errorf("error creating a new user: %w", err)
	}

	u.ID = id
	return u, nil
}

func (uc *userUsercase) List(ctx context.Context, limit, offset int) ([]*model.User, error) {
	// TODO - validate api key and admin user permission

	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	users, err := uc.store.List(ctx, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("error listing users: %w", err)
	}

	return users, nil
}

func (uc *userUsercase) Get(ctx context.Context, id int) (*model.User, error) {
	// TODO - validate api key and admin user permission

	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	u, err := uc.store.Get(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("error fetching user: %w", err)
	}

	return u, nil
}

func (uc *userUsercase) Update(ctx context.Context, u *model.User) (*model.User, error) {
	// TODO - validate ApiKey, base user only can update their data, admin full rw

	_, err := uc.store.Get(ctx, u.ID)
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
	// TODO - validate apikey and permissions base user can delete their account admin full rw

	if err := uc.store.Delete(ctx, id); err != nil {
		return fmt.Errorf("error deleting user: %w", err)
	}

	return nil
}
