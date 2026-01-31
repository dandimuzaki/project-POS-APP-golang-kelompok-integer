package repository

import (
	"context"
	"project-POS-APP-golang-integer/internal/data/entity"
	"project-POS-APP-golang-integer/internal/infra"
	"time"

	"go.uber.org/zap"
	"gorm.io/gorm"
)

type PasswordResetRepository interface {
	Create(ctx context.Context, reset *entity.PasswordReset) error
	MarkUsed(ctx context.Context, tokenID uint) error
	GetValidByUser(ctx context.Context, userID uint) (*entity.PasswordReset, error)
}

type passwordResetRepository struct {
	db     *gorm.DB
	Logger *zap.Logger
}

func NewPasswordResetRepo(db *gorm.DB, log *zap.Logger) PasswordResetRepository {
	return &passwordResetRepository{
		db:     db,
		Logger: log,
	}
}

func (r *passwordResetRepository) Create(ctx context.Context, reset *entity.PasswordReset) error {
	db := infra.GetDB(ctx, r.db)
	if err := db.Create(reset).Error; err != nil {
		r.Logger.Error("Error create reset password", zap.Error(err))
		return err
	}
	return nil
}

func (r *passwordResetRepository) GetValidByUser(ctx context.Context, userID uint) (*entity.PasswordReset, error) {
	var reset entity.PasswordReset

	err := r.db.
		Where(
			"user_id = ? AND expired_at > ? AND used_at IS NULL",
			userID,
			time.Now(),
		).
		Order("created_at DESC").
		First(&reset).Error

	if err != nil {
		r.Logger.Error("Error get valid reset password", zap.Error(err))
		return nil, err
	}

	return &reset, nil
}

func (r *passwordResetRepository) MarkUsed(ctx context.Context, tokenID uint) error {
	err := r.db.
		Model(&entity.PasswordReset{}).
		Where("id = ?", tokenID).
		Update("used_at", time.Now()).
		Error

	if err != nil {
		r.Logger.Error("Error mark reset password used", zap.Error(err))
	}

	return err
}
