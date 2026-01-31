package usecase

import (
	"context"
	"errors"
	"project-POS-APP-golang-integer/internal/data/entity"
	"project-POS-APP-golang-integer/internal/data/repository"
	"project-POS-APP-golang-integer/internal/data/repository/mocks"
	"project-POS-APP-golang-integer/internal/dto/request"
	"project-POS-APP-golang-integer/internal/infra"
	"project-POS-APP-golang-integer/pkg/utils"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.uber.org/zap"
)

func TestUserService_GetUserByID_Success(t *testing.T) {
	ctx := context.Background()

	userRepo := new(mocks.UserRepoMock)
	profileRepo := new(mocks.ProfileRepoMock)
	tx := new(infra.MockTxManager)

	repo := repository.Repository{
		UserRepo: userRepo,
		ProfileRepo: profileRepo,
	}

	log := zap.NewNop()
	service := NewUserService(tx, &repo, log)

	expectedUser := entity.User{
		ID:    1,
		Email: "test@mail.com",
		Role:  entity.RoleAdmin,
	}

	userRepo.
		On("GetUserByID", ctx, uint(1)).
		Return(expectedUser, nil)

	user, err := service.GetUserByID(ctx, 1)

	assert.NoError(t, err)
	assert.Equal(t, expectedUser.Email, user.Email)
	assert.Equal(t, expectedUser.Role, user.Role)
	userRepo.AssertExpectations(t)
}

func TestUserService_GetUserByID_Error(t *testing.T) {
	ctx := context.Background()

	userRepo := new(mocks.UserRepoMock)
	profileRepo := new(mocks.ProfileRepoMock)
	tx := new(infra.MockTxManager)

	repo := repository.Repository{
		UserRepo:    userRepo,
		ProfileRepo: profileRepo,
	}

	log := zap.NewNop()
	service := NewUserService(tx, &repo, log)

	expectedErr := utils.ErrUserNotFound
	userRepo.
		On("GetUserByID", ctx, uint(2)).
		Return(entity.User{}, expectedErr)

	user, err := service.GetUserByID(ctx, 2)

	assert.Error(t, err)
	assert.Equal(t, expectedErr, err)
	assert.Equal(t, entity.User{}, user)

	userRepo.AssertExpectations(t)
}

func TestUserService_CreateUser_Success(t *testing.T) {
	ctx := context.Background()

	userRepo := new(mocks.UserRepoMock)
	profileRepo := new(mocks.ProfileRepoMock)
	tx := new(infra.MockTxManager)

	repo := repository.Repository{
		UserRepo:    userRepo,
		ProfileRepo: profileRepo,
	}

	log := zap.NewNop()
	service := NewUserService(tx, &repo, log)

	req := request.UserRequest{
		Email: "test@mail.com",
		Role:  "admin",
	}

	createdUser := &entity.User{
		ID:    1,
		Email: req.Email,
		Role:  entity.RoleAdmin,
	}

	tx.
		On("WithinTx", mock.Anything).
		Return(nil)

	userRepo.
		On("CreateUser", mock.Anything, mock.AnythingOfType("*entity.User")).
		Return(createdUser, nil)

	profileRepo.
		On("CreateProfile", mock.Anything, mock.AnythingOfType("*entity.Profile")).
		Return(&entity.Profile{}, nil)

	res, err := service.CreateUser(ctx, req)

	assert.NoError(t, err)
	assert.Equal(t, req.Email, res.Email)
	assert.Equal(t, entity.RoleAdmin, res.Role)
	assert.NotEmpty(t, res.Password) // generated password

	userRepo.AssertExpectations(t)
	profileRepo.AssertExpectations(t)
	tx.AssertExpectations(t)
}

func TestUserService_CreateUser_UserRepoError(t *testing.T) {
	ctx := context.Background()

	userRepo := new(mocks.UserRepoMock)
	profileRepo := new(mocks.ProfileRepoMock)
	tx := new(infra.MockTxManager)

	repo := repository.Repository{
		UserRepo:    userRepo,
		ProfileRepo: profileRepo,
	}

	log := zap.NewNop()
	service := NewUserService(tx, &repo, log)

	req := request.UserRequest{
		Email: "test@mail.com",
		Role:  "admin",
	}

	expectedErr := errors.New("db error")

	// Tx always runs
	tx.On("WithinTx", mock.Anything).Return(nil)

	// User creation fails
	userRepo.
		On("CreateUser", mock.Anything, mock.AnythingOfType("*entity.User")).
		Return((*entity.User)(nil), expectedErr)

	res, err := service.CreateUser(ctx, req)

	assert.Error(t, err)
	assert.Nil(t, res)
	assert.Equal(t, expectedErr, err)

	userRepo.AssertExpectations(t)
	profileRepo.AssertNotCalled(t, "CreateProfile")
}

func TestUserService_CreateUser_ProfileRepoError(t *testing.T) {
	ctx := context.Background()

	userRepo := new(mocks.UserRepoMock)
	profileRepo := new(mocks.ProfileRepoMock)
	tx := new(infra.MockTxManager)

	repo := repository.Repository{
		UserRepo:    userRepo,
		ProfileRepo: profileRepo,
	}

	log := zap.NewNop()
	service := NewUserService(tx, &repo, log)

	req := request.UserRequest{
		Email: "test@mail.com",
		Role:  "admin",
	}

	tx.On("WithinTx", mock.Anything).Return(nil)

	createdUser := &entity.User{
		ID:    1,
		Email: req.Email,
		Role:  entity.RoleAdmin,
		PasswordHash: "123",
	}

	userRepo.
		On("CreateUser", mock.Anything, mock.AnythingOfType("*entity.User")).
		Return(createdUser, nil)

	expectedErr := errors.New("profile create failed")

	profileRepo.
		On("CreateProfile", mock.Anything, mock.AnythingOfType("*entity.Profile")).
		Return((*entity.Profile)(nil), expectedErr)

	res, err := service.CreateUser(ctx, req)

	assert.Error(t, err)
	assert.Nil(t, res)
	assert.Equal(t, expectedErr, err)

	userRepo.AssertExpectations(t)
	profileRepo.AssertExpectations(t)
}
