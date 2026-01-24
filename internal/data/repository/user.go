package repository

import (
	"context"
	"project-POS-APP-golang-integer/internal/data/entity"
	"project-POS-APP-golang-integer/internal/dto"

	"go.uber.org/zap"
	"gorm.io/gorm"
)

type UserRepository interface{
	GetListUsers(ctx context.Context, req dto.UserFilterRequest) ([]entity.User, int, error)
}

type userRepository struct {
	db *gorm.DB
	Logger *zap.Logger
}

func NewUserRepo(db *gorm.DB, log *zap.Logger) UserRepository {
	return &userRepository{
		db: db,
		Logger: log,
	}
}

func (r *userRepository) GetListUsers(ctx context.Context, req dto.UserFilterRequest) ([]entity.User, int, error) {
	return nil, 0, nil
}