package usecase

import (
	"context"
	"errors"
	"project-POS-APP-golang-integer/internal/data/entity"
	"project-POS-APP-golang-integer/internal/data/repository"
	"project-POS-APP-golang-integer/internal/data/repository/mocks"
	"project-POS-APP-golang-integer/internal/dto/request"
	"testing"
	"time"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

func TestProfileService_GetProfile_Success(t *testing.T) {
	require := require.New(t)

	mockRepo := new(mocks.ProfileRepoMock)

	repo := &repository.Repository{
		ProfileRepo: mockRepo,
	}

	service := NewProfileService(nil, repo, zap.NewNop())

	ctx := context.WithValue(context.Background(), "user_id", uint(1))

	mockUser := &entity.User{
		Email: "user@test.com",
		Role:  "admin",
		Profile: &entity.Profile{
			FullName:    "Dandi",
			Phone:       "08123",
			DateOfBirth: time.Date(2000, 1, 2, 0, 0, 0, 0, time.UTC),
			Salary:      5000,
		},
	}

	mockRepo.
		On("GetProfileByID", ctx, uint(1)).
		Return(mockUser, nil).
		Once()

	result, err := service.GetProfile(ctx)

	require.NoError(err)
	require.NotNil(result)
	require.Equal("user@test.com", result.Email)
	require.Equal("Dandi", result.FullName)

	mockRepo.AssertExpectations(t)
}

func TestProfileService_GetProfile_Error(t *testing.T) {
	require := require.New(t)

	mockRepo := new(mocks.ProfileRepoMock)

	repo := &repository.Repository{
		ProfileRepo: mockRepo,
	}

	service := NewProfileService(nil, repo, zap.NewNop())

	ctx := context.WithValue(context.Background(), "user_id", uint(1))

	mockRepo.
		On("GetProfileByID", ctx, uint(1)).
		Return((*entity.User)(nil), errors.New("db error")).
		Once()

	result, err := service.GetProfile(ctx)

	require.Error(err)
	require.Nil(result)

	mockRepo.AssertExpectations(t)
}

func TestProfileService_UpdateProfile_Success(t *testing.T) {
	require := require.New(t)

	mockRepo := new(mocks.ProfileRepoMock)

	repo := &repository.Repository{
		ProfileRepo: mockRepo,
	}

	service := NewProfileService(nil, repo, zap.NewNop())

	ctx := context.WithValue(context.Background(), "user_id", uint(1))

	req := &request.ProfileRequest{
		FullName: "Dandi",
		DateOfBirth: "2 January 2000",
	}

	mockRepo.
		On("UpdateProfile", ctx, mock.AnythingOfType("*entity.Profile")).
		Return(nil).
		Once()

	err := service.UpdateProfile(ctx, req)

	require.NoError(err)

	mockRepo.AssertExpectations(t)
}

func TestProfileService_UpdateProfile_Error(t *testing.T) {
	require := require.New(t)

	mockRepo := new(mocks.ProfileRepoMock)

	repo := &repository.Repository{
		ProfileRepo: mockRepo,
	}

	service := NewProfileService(nil, repo, zap.NewNop())

	ctx := context.WithValue(context.Background(), "user_id", uint(1))

	req := &request.ProfileRequest{
		FullName: "Fail",
		DateOfBirth: "2 January 2000",
	}

	mockRepo.
		On("UpdateProfile", ctx, mock.Anything).
		Return(errors.New("update failed")).
		Once()

	err := service.UpdateProfile(ctx, req)

	require.Error(err)

	mockRepo.AssertExpectations(t)
}
