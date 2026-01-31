package repository

import (
	"context"
	"project-POS-APP-golang-integer/internal/data/entity"
	"project-POS-APP-golang-integer/internal/infra"
	"time"

	"go.uber.org/zap"
	"gorm.io/gorm"
)

type OTPRepository interface {
	Create(ctx context.Context, otp *entity.OTP) error
	MarkUsed(ctx context.Context, otpID uint) error
	GetValidByUser(ctx context.Context, userID uint) (*entity.OTP, error)
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

func (r *otpRepository) Create(ctx context.Context, otp *entity.OTP) error {
	db := infra.GetDB(ctx, r.db)
	if err := db.Create(otp).Error; err != nil {
		r.Logger.Error("Error create OTP", zap.Error(err))
		return err
	}
	return nil
}

func (r *otpRepository) GetValidByUser(ctx context.Context, userID uint) (*entity.OTP, error) {
	var otp entity.OTP

	err := r.db.
		Where(
			"user_id = ? AND expires_at > ? AND is_used = false",
			userID,
			time.Now(),
		).
		Order("created_at DESC").
		First(&otp).Error

	if err != nil {
		r.Logger.Error("Error get valid OTP", zap.Error(err))
		return nil, err
	}

	return &otp, nil
}

func (r *otpRepository) MarkUsed(ctx context.Context, otpID uint) error {
	err := r.db.
		Model(&entity.OTP{}).
		Where("id = ?", otpID).
		Update("is_used", true).
		Error

	if err != nil {
		r.Logger.Error("Error mark OTP used", zap.Error(err))
	}

	return err
}
