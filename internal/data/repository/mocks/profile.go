package mocks

import (
	"context"

	"project-POS-APP-golang-integer/internal/data/entity"

	"github.com/stretchr/testify/mock"
)

type ProfileRepoMock struct {
	mock.Mock
}

func (m *ProfileRepoMock) GetProfileByID(ctx context.Context, id uint) (*entity.User, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(*entity.User), args.Error(1)
}

func (m *ProfileRepoMock) CreateProfile(ctx context.Context, profile *entity.Profile) (*entity.Profile, error) {
	args := m.Called(ctx, profile)
	return args.Get(0).(*entity.Profile), args.Error(1)
}

func (m *ProfileRepoMock) UpdateProfile(ctx context.Context, profile *entity.Profile) error {
	args := m.Called(ctx, profile)
	return args.Error(0)
}
