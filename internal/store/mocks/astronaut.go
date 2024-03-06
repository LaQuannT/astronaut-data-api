package mocks

import (
	"context"

	"github.com/LaQuannT/astronaut-data-api/internal/model"
	"github.com/stretchr/testify/mock"
)

type AstronautStore struct {
	mock.Mock
}

func (m *AstronautStore) Create(ctx context.Context, a *model.Astronaut) (int, error) {
	args := m.Called(ctx, a)
	return args.Int(0), args.Error(1)
}

func (m *AstronautStore) List(ctx context.Context, limit, offset int) ([]*model.Astronaut, error) {
	args := m.Called(ctx, limit, offset)
	return args.Get(0).([]*model.Astronaut), args.Error(1)
}

func (m *AstronautStore) Get(ctx context.Context, id int) (*model.Astronaut, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(*model.Astronaut), args.Error(1)
}

func (m *AstronautStore) Update(ctx context.Context, a *model.Astronaut) error {
	args := m.Called(ctx, a)
	return args.Error(0)
}

func (m *AstronautStore) Delete(ctx context.Context, id int) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}
