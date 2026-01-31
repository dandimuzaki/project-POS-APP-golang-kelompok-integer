package usecase

import (
	"context"
	"errors"
	"project-POS-APP-golang-integer/internal/data/entity"
	"project-POS-APP-golang-integer/internal/data/repository"
	"project-POS-APP-golang-integer/internal/dto/request"
	"project-POS-APP-golang-integer/internal/dto/response"
	"project-POS-APP-golang-integer/pkg/utils"
	content "project-POS-APP-golang-integer/pkg/utils/email"
	"time"

	"go.uber.org/zap"
	"gorm.io/gorm"
)

type AuthService interface {
	Login(ctx context.Context, data request.LoginRequest) (*response.AuthResponse, error)
	ValidateToken(ctx context.Context, token string) (*uint, error)
	Logout(ctx context.Context, token string) error
	RequestResetPassword(ctx context.Context, email string) (*response.OTPResponse, error)
	ValidateOTP(ctx context.Context, req request.ValidateOTP) (*response.ResetTokenResponse, error)
	ResetPassword(ctx context.Context, req request.ResetPassword) error
}

type authService struct {
	tx TxManager
	repo *repository.Repository
	log *zap.Logger
	email EmailSender
}

func NewAuthService(tx TxManager, repo *repository.Repository, log *zap.Logger, email EmailSender) AuthService {
	return &authService{
		tx: tx,
		repo: repo,
		log: log,
		email: email,
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
	err = s.repo.OTPRepo.Create(ctx, &otp)
	if err != nil {
		s.log.Error("Error create otp", zap.Error(err))
		return nil, err
	}

	res := response.OTPResponse{
		OTPCode: otpCode,
		ExpiresAt: otp.ExpiresAt,
	}

	// Send email
	err = s.email.Send(ctx, request.EmailRequest{
		To:      user.Email,
		Subject: "Reset Password",
		Body:    content.SendOTP(res),
	})

	if err != nil {
		s.log.Error("Error send email", zap.Error(err))
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
	otp, err := s.repo.OTPRepo.GetValidByUser(ctx, user.ID)
	if err != nil {
		s.log.Error("Error get valid OTP", zap.Error(err))
		return nil, err
	}

	if otp.OTPCode != req.OTP {
		return nil, err
	}

	// Mark OTP as used
	err = s.repo.OTPRepo.MarkUsed(ctx, otp.ID)
	if err != nil {
		s.log.Error("Error mark OTP as used", zap.Error(err))
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
	err = s.repo.PasswordResetRepo.Create(ctx, &reset)
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
	// Get user by email
	user, err := s.repo.UserRepo.FindUserByEmail(ctx, req.Email)
	if err != nil {
		s.log.Error("Error get user by id", zap.Error(err))
		return err
	}

	// Get valid token
	reset, err := s.repo.PasswordResetRepo.GetValidByUser(ctx, user.ID)
	if err != nil {
		s.log.Error("Error get token", zap.Error(err))
		return err
	}

	if !utils.CheckPassword(req.ResetToken, reset.ResetTokenHash) {
		return utils.ErrInvalidToken
	}

	err = s.tx.WithinTx(ctx, func(ctx context.Context) error {
		// Hash password
		user.PasswordHash = utils.HashPassword(req.NewPassword)
		// Update user password
		e := s.repo.UserRepo.UpdateUser(ctx, user.ID, user)
		if e != nil {
			return e
		}
		// Mark token as used
		e = s.repo.PasswordResetRepo.MarkUsed(ctx, reset.ID)
		if e != nil {
			return e
		}
		return nil
	})

	if err != nil {
		s.log.Error("Error reset password", zap.Error(err))
		return err
	}

	return nil
}