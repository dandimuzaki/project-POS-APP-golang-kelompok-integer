package repository

import (
	"context"
	"project-POS-APP-golang-integer/internal/data/entity"
	"project-POS-APP-golang-integer/internal/infra"
	"project-POS-APP-golang-integer/pkg/utils"

	"go.uber.org/zap"
	"gorm.io/gorm"
)

type PasswordResetRepository interface {
	CreateResetToken(ctx context.Context, reset entity.PasswordReset) error
	MarkResetToken(ctx context.Context, token string) error
	ValidateResetToken(ctx context.Context, token string) (*uint, error)
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

func (r *passwordResetRepository) CreateResetToken(ctx context.Context, reset entity.PasswordReset) error {
	db := infra.GetDB(ctx, r.db)
	err := db.Create(&reset).Error
	if err != nil {
		r.Logger.Error("Error query create token: ", zap.Error(err))
		return err
	}

	return nil
}

func (r *passwordResetRepository) ValidateResetToken(ctx context.Context, token string) (*uint, error) {
	db := infra.GetDB(ctx, r.db)
	// Validate token to authorize user
	type result struct {
		UserID uint
	}

	var res result
	tokenHash := utils.HashPassword(token)

	err := db.
		Model(&entity.PasswordReset{}).
		Select("user_id").
		Where("reset_token_hash = ? AND expired_at > NOW() AND used_at IS NULL", tokenHash).
		First(&res).
		Error
	if err != nil {
		r.Logger.Error("Error query validate token: ", zap.Error(err))
		return nil, err
	}

	return &res.UserID, nil
}

func (r *passwordResetRepository) MarkResetToken(ctx context.Context, token string) error {
	// Mark token after request success
	err := r.db.Model(&entity.PasswordReset{}).Where("reset_token_hash = ?", token).Update("used_at", "NOW()").Error
	if err != nil {
		r.Logger.Error("Error query update token: ", zap.Error(err))
	}

	return err
}