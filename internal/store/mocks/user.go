package mocks

import (
	"context"

	"github.com/LaQuannT/astronaut-data-api/internal/model"
	"github.com/stretchr/testify/mock"
)

type UserStore struct {
	mock.Mock
}

func (m *UserStore) Create(ctx context.Context, u *model.User) (int, error) {
	args := m.Called(ctx, u)
	return args.Int(0), args.Error(1)
}

func (m *UserStore) List(ctx context.Context, limit, offset int) ([]*model.User, error) {
	args := m.Called(ctx, limit, offset)
	return args.Get(0).([]*model.User), args.Error(1)
}

func (m *UserStore) Get(ctx context.Context, id int) (*model.User, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(*model.User), args.Error(1)
}

func (m *UserStore) Update(ctx context.Context, u *model.User) error {
	args := m.Called(ctx, u)
	return args.Error(0)
}

func (m *UserStore) Delete(ctx context.Context, id int) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *UserStore) SearchApiKey(ctx context.Context, key string) (*model.User, error) {
	args := m.Called(ctx, key)
	return args.Get(0).(*model.User), args.Error(1)
}
