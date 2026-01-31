package usecase

import (
	"context"
	"errors"
	"project-POS-APP-golang-integer/internal/data/entity"
	"project-POS-APP-golang-integer/internal/data/repository"
	"project-POS-APP-golang-integer/internal/dto/request"
	"project-POS-APP-golang-integer/internal/dto/response"
	"project-POS-APP-golang-integer/pkg/utils"
	"time"

	"go.uber.org/zap"
	"gorm.io/gorm"
)

type AuthService interface {
	Login(ctx context.Context, data request.LoginRequest) (*response.AuthResponse, error)
	ValidateToken(ctx context.Context, token string) (*uint, error)
	Logout(ctx context.Context, token string) error
	RequestResetPassword(ctx context.Context, email string) (*response.OTPResponse, error)
	ResetPassword(ctx context.Context, req request.ResetPassword) error
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

func (s *authService) Login(ctx context.Context, data request.LoginRequest) (*response.AuthResponse, error) {
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
	
	res := response.AuthResponse{
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

func (s *authService) RequestResetPassword(ctx context.Context, email string) (*response.OTPResponse, error) {
	// Find user by email
	user, err := s.repo.UserRepo.FindUserByEmail(ctx, email)
	if err != nil {
		s.log.Error("Error find user by email", zap.Error(err))
		return nil, err
	}

	// Generate OTP
	otpCode, _ := utils.GenerateOTP(6)
	otp := entity.OTP{
		UserID: user.ID,
		OTPCode: otpCode,
		ExpiresAt: time.Now().Add(10 * time.Minute),
	}
	err = s.repo.OTPRepo.CreateOTP(ctx, otp)
	if err != nil {
		s.log.Error("Error create otp", zap.Error(err))
		return nil, err
	}

	// Send email with OTP

	res := response.OTPResponse{
		OTPCode: otpCode,
		ExpiresAt: otp.ExpiresAt,
	}
	
	return &res, nil
}

func (s *authService) ValidateOTP(ctx context.Context, req request.ValidateOTP) (*response.ResetTokenResponse, error) {
	// Find user by email
	user, err := s.repo.UserRepo.FindUserByEmail(ctx, req.Email)
	if err != nil {
		s.log.Error("Error find user by email", zap.Error(err))
		return nil, err
	}

	// Validate OTP
	err = s.repo.OTPRepo.ValidateOTP(ctx, req.OTP)
	if err != nil {
		s.log.Error("Error validate OTP", zap.Error(err))
		return nil, err
	}

	// Generate reset token
	token, err := utils.GenerateRandomToken(16)
	reset := entity.PasswordReset{
		UserID: user.ID,
		ResetTokenHash: utils.HashPassword(token.String()),
		ExpiredAt: time.Now().Add(5 * time.Minute),
		CreatedAt: time.Now(),
	}
	err = s.repo.PasswordResetRepo.CreateResetToken(ctx, reset)
	if err != nil {
		s.log.Error("Error create reset token", zap.Error(err))
		return nil, err
	}

	res := response.ResetTokenResponse{
		ResetToken: token.String(),
	}

	return &res, nil
}

func (s *authService) ResetPassword(ctx context.Context, req request.ResetPassword) error {
	userID, err := s.repo.PasswordResetRepo.ValidateResetToken(ctx, req.ResetToken)
	if err != nil {
		s.log.Error("Error validate token", zap.Error(err))
		return err
	}

	// Get user by id
	user, err := s.repo.UserRepo.GetUserByID(ctx, *userID)
	if err != nil {
		s.log.Error("Error get user by id", zap.Error(err))
		return err
	}	

	// Hash password
	user.PasswordHash = utils.HashPassword(req.NewPassword)

	err = s.tx.WithinTx(ctx, func(ctx context.Context) error {
		// Update user password
		e := s.repo.UserRepo.UpdateUser(ctx, user.ID, &user)
		if e != nil {
			return e
		}
		// Mark OTP as used
		e = s.repo.OTPRepo.MarkOTP(ctx, req.OTP)
		if e != nil {
			return e
		}
		return nil
	})

	return nil
}