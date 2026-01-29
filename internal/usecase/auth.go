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
	tx TxManager
	repo *repository.Repository
	log *zap.Logger
	config utils.Configuration
}

func NewAuthService(tx TxManager, repo *repository.Repository, log *zap.Logger, config utils.Configuration) AuthService {
	return &authService{
		tx: tx,
		repo: repo,
		log: log,
		config: config,
	}
}

func (s *authService) Login(ctx context.Context, data dto.LoginRequest) (*dto.AuthResponse, error) {
	// Find user by email
	user, err := s.repo.UserRepo.FindUserByEmail(ctx, data.Email)
	if err == gorm.ErrRecordNotFound {
		s.log.Error("User not found: ", zap.Error(err))
		return nil, err
	}

	if err != nil {
		s.log.Error("Error find user by email: ", zap.Error(err))
		return nil, err
	}

	// Check password
	if !utils.CheckPassword(data.Password, user.PasswordHash) {
		s.log.Error("Incorrect password: ", zap.Error(err))
		return nil, errors.New("incorrect password")
	}

	// Record session
	token, err := s.repo.SessionRepo.Create(ctx, user.ID)
	if err != nil {
		s.log.Error("Error create token: ", zap.Error(err))
		return nil, errors.New("token error")
	}
	
	res := dto.AuthResponse{
		Token: token,
	}

	return &res, nil
}

func (s *authService) ValidateToken(ctx context.Context, token string) (*uint, error) {
	// Validate token to authorize user
	userID, err := s.repo.SessionRepo.ValidateToken(ctx, token)
	if err != nil {
		s.log.Error("Error validate token service: ", zap.Error(err))
		return nil, err
	}

	return userID, nil
}

func (s *authService) Logout(ctx context.Context, token string) error {
	err := s.repo.SessionRepo.Revoke(ctx, token)
	if err != nil {
		s.log.Error("Error logout service: ", zap.Error(err))
		return err
	}
	return nil
}