package usecase

import (
	"context"
	"errors"
	"project-POS-APP-golang-integer/internal/data/repository"
	"project-POS-APP-golang-integer/internal/dto"
	"project-POS-APP-golang-integer/pkg/utils"

	"go.uber.org/zap"
	"gorm.io/gorm"
)

type AuthService interface {
	Login(ctx context.Context, data dto.LoginRequest) (*dto.AuthResponse, error)
	ValidateToken(ctx context.Context, token string) (*uint, error)
	Logout(ctx context.Context, token string) error
}

type authService struct {
	Repo *repository.Repository
	Logger *zap.Logger
	Config utils.Configuration
}

func NewAuthService(repo *repository.Repository, log *zap.Logger, config utils.Configuration) AuthService {
	return &authService{
		Repo: repo,
		Logger: log,
		Config: config,
	}
}

func (u *authService) Login(ctx context.Context, data dto.LoginRequest) (*dto.AuthResponse, error) {
	// Find user by email
	user, err := u.Repo.UserRepo.FindByEmail(ctx, data.Email)
	if err == gorm.ErrRecordNotFound {
		u.Logger.Error("User not found: ", zap.Error(err))
		return nil, err
	}

	if err != nil {
		u.Logger.Error("Error find user by email: ", zap.Error(err))
		return nil, err
	}

	// Check password
	if !utils.CheckPassword(data.Password, user.PasswordHash) {
		u.Logger.Error("Incorrect password: ", zap.Error(err))
		return nil, errors.New("incorrect password")
	}

	// Record session
	token, err := u.Repo.SessionRepo.Create(ctx, user.ID)
	if err != nil {
		u.Logger.Error("Error create token: ", zap.Error(err))
		return nil, errors.New("token error")
	}
	
	res := dto.AuthResponse{
		Token: token,
	}

	return &res, nil
}

func (u *authService) ValidateToken(ctx context.Context, token string) (*uint, error) {
	// Validate token to authorize user
	userID, err := u.Repo.SessionRepo.ValidateToken(ctx, token)
	if err != nil {
		u.Logger.Error("Error validate token service: ", zap.Error(err))
		return nil, err
	}

	return userID, nil
}

func (u *authService) Logout(ctx context.Context, token string) error {
	err := u.Repo.SessionRepo.Revoke(ctx, token)
	if err != nil {
		u.Logger.Error("Error logout service: ", zap.Error(err))
		return err
	}
	return nil
}