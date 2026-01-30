package repository

import (
	"context"
	"project-POS-APP-golang-integer/internal/data/entity"
	"project-POS-APP-golang-integer/internal/infra"

	"go.uber.org/zap"
	"gorm.io/gorm"
)

type OTPRepository interface {
	CreateOTP(ctx context.Context, otp entity.OTP) error
	MarkOTP(ctx context.Context, token string) error
	ValidateOTP(ctx context.Context, token string) error
}

type otpRepository struct {
	db     *gorm.DB
	Logger *zap.Logger
}

func NewOTPRepo(db *gorm.DB, log *zap.Logger) OTPRepository {
	return &otpRepository{
		db:     db,
		Logger: log,
	}
}

func (r *otpRepository) CreateOTP(ctx context.Context, otp entity.OTP) error {
	db := infra.GetDB(ctx, r.db)
	err := db.Create(&otp).Error
	if err != nil {
		r.Logger.Error("Error query create OTP: ", zap.Error(err))
		return err
	}

	return nil
}

func (r *otpRepository) ValidateOTP(ctx context.Context, otpCode string) error {
	var otp entity.OTP
	// Validate otp to proceed request
	query := r.db.Model(&otp).Where("otp_code = ? AND expires_at > NOW() AND is_used IS false", otp).Limit(1)
	err := query.Find(&otp).Error
	if err != nil {
		r.Logger.Error("Error query validate otp: ", zap.Error(err))
		return err
	}

	return nil
}

func (r *otpRepository) MarkOTP(ctx context.Context, otp string) error {
	// Mark otp after request success
	err := r.db.Model(&entity.OTP{}).Where("otp_code = ?", otp).Update("is_used", "true").Error
	if err != nil {
		r.Logger.Error("Error query update OTP: ", zap.Error(err))
	}

	return err
}