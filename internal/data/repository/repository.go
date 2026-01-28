package repository

import (
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type Repository struct {
	UserRepo UserRepository
	SessionRepo SessionRepository
}

func NewRepository(db *gorm.DB, log *zap.Logger) *Repository {
	return &Repository{
		UserRepo: NewUserRepo(db, log),
		SessionRepo: NewSessionRepo(db, log),
	}
}
