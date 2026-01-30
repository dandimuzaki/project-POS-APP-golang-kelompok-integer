package repository

import (
	"context"
	"project-POS-APP-golang-integer/internal/data/entity"
	"project-POS-APP-golang-integer/internal/infra"
	"project-POS-APP-golang-integer/pkg/utils"
	"time"

	"github.com/google/uuid"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type SessionRepository interface {
	Create(ctx context.Context, userID uint) (uuid.UUID, error)
	Revoke(ctx context.Context, token string) error
	ValidateToken(ctx context.Context, token string) (*uint, error)
}

type sessionRepository struct {
	db     *gorm.DB
	Logger *zap.Logger
}

func NewSessionRepo(db *gorm.DB, log *zap.Logger) SessionRepository {
	return &sessionRepository{
		db:     db,
		Logger: log,
	}
}

func (r *sessionRepository) Create(ctx context.Context, userID uint) (uuid.UUID, error) {
	db := infra.GetDB(ctx, r.db)
	// Create session after login and register
	token, err := utils.GenerateRandomToken(16)
	if err != nil {
		r.Logger.Error("Error create token: ", zap.Error(err))
		return uuid.Nil, err
	}

	session := entity.Session{
		UserID:    uint(userID),
		Token:     token,
		ExpiresAt: time.Now().AddDate(0, 0, 5),
		CreatedAt: time.Now(),
	}

	err = db.Create(&session).Error
	if err != nil {
		r.Logger.Error("Error query create session: ", zap.Error(err))
		return uuid.Nil, err
	}

	return token, nil
}

func (r *sessionRepository) ValidateToken(ctx context.Context, token string) (*uint, error) {
	db := infra.GetDB(ctx, r.db)
	// Validate token to authorize user
	type result struct {
		UserID uint
	}

	var res result

	err := db.
		Model(&entity.Session{}).
		Select("user_id").
		Where("token = ? AND expires_at > NOW() AND revoked_at IS NULL", token).
		First(&res).
		Error
	if err != nil {
		r.Logger.Error("Error query validate token: ", zap.Error(err))
		return nil, err
	}

	return &res.UserID, nil
}

func (r *sessionRepository) Revoke(ctx context.Context, token string) error {
	db := infra.GetDB(ctx, r.db)
	// Revoke session after logout
	err := db.Model(&entity.Session{}).Where("token = ?", token).Update("revoked_at", "NOW()").Error
	if err != nil {
		r.Logger.Error("Error query revoke session: ", zap.Error(err))
	}

	return err
}