package mocks

import (
	"context"

	"project-POS-APP-golang-integer/internal/data/entity"
	"project-POS-APP-golang-integer/internal/data/repository"

	"github.com/stretchr/testify/mock"
)

type UserRepoMock struct {
	mock.Mock
}

func (m *UserRepoMock) GetUserList(ctx context.Context, f repository.UserQueryParams) ([]entity.User, int64, error) {
	args := m.Called(ctx, f)
	return args.Get(0).([]entity.User), args.Get(1).(int64), args.Error(2)
}

func (m *UserRepoMock) FindUserByEmail(ctx context.Context, email string) (*entity.User, error) {
	args := m.Called(ctx, email)
	return args.Get(0).(*entity.User), args.Error(1)
}

func (m *UserRepoMock) GetUserByID(ctx context.Context, id uint) (entity.User, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(entity.User), args.Error(1)
}

func (m *UserRepoMock) CreateUser(ctx context.Context, user *entity.User) (*entity.User, error) {
	args := m.Called(ctx, user)
	return args.Get(0).(*entity.User), args.Error(1)
}

func (m *UserRepoMock) UpdateUser(ctx context.Context, id uint, user *entity.User) error {
	args := m.Called(ctx, id, user)
	return args.Error(0)
}

func (m *UserRepoMock) DeleteUser(ctx context.Context, id uint) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}