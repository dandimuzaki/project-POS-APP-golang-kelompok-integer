package repository

import (
	"context"
	"project-POS-APP-golang-integer/internal/data/entity"
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

	tx := r.db.Begin()
	if tx.Error != nil {
		return uuid.Nil, tx.Error
	}
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		} else {
			tx.Commit()
		}
	}()

	err = tx.Create(&session).Error
	if err != nil {
		r.Logger.Error("Error query create session: ", zap.Error(err))
		return uuid.Nil, err
	}

	return token, nil
}

func (r *sessionRepository) ValidateToken(ctx context.Context, token string) (*uint, error) {
	// Validate token to authorize user
	var userID *uint
	query := r.db.Model(&entity.Session{}).Select("user_id").Where("token = ?", token).Where("expired_at > NOW()").Where("revoked_at IS NULL")
	err := query.Find(&userID).Error
	if err != nil {
		r.Logger.Error("Error query validate token: ", zap.Error(err))
		return nil, err
	}

	return userID, nil
}

func (r *sessionRepository) Revoke(ctx context.Context, token string) error {
	// Revoke session after logout
	err := r.db.Model(&entity.Session{}).Where("token = ?", token).Update("revoked_at", "NOW()").Error
	if err != nil {
		r.Logger.Error("Error query revoke session: ", zap.Error(err))
	}

	return err
}