package repository

import (
	"context"
	"project-POS-APP-golang-integer/internal/data/entity"
	"project-POS-APP-golang-integer/internal/dto"
	"project-POS-APP-golang-integer/internal/infra"

	"go.uber.org/zap"
	"gorm.io/gorm"
)

type UserRepository interface {
	CreateUser(ctx context.Context, user *entity.User) (*entity.User, error)
	FindUserByEmail(ctx context.Context, email string) (*entity.User, error)
	GetUserList(ctx context.Context, f dto.UserFilterRequest) ([]entity.User, int64, error)
	GetUserByID(ctx context.Context, id uint) (entity.User, error)
	UpdateUser(ctx context.Context, id uint, data *entity.User) error
	DeleteUser(ctx context.Context, id uint) error
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

func (r *userRepository) CreateUser(ctx context.Context, user *entity.User) (*entity.User, error) {
	db := infra.GetDB(ctx, r.db)
	err := db.Create(&user).Error
	if err != nil {
		r.Logger.Error("Error query create user", zap.Error(err))
		return nil, err
	}
	return user, err
}

func (r *userRepository) FindUserByEmail(ctx context.Context, email string) (*entity.User, error) {
	db := infra.GetDB(ctx, r.db)
	var user entity.User
	query := db.WithContext(ctx).Model(&user).Where("email = ?", email).Limit(1)
	err := query.Find(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, err
}

func (r *userRepository) GetUserList(ctx context.Context, f dto.UserFilterRequest) ([]entity.User, int64, error) {
	db := infra.GetDB(ctx, r.db)
	var users []entity.User
	var total int64
	query := db.Model(&entity.User{})

	// Filter by role
	query = query.Where("role = ?", f.Role)

	// Get total user
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (f.Page - 1) * f.Limit
	err := query.Limit(f.Limit).Offset(offset).Find(&users).Error
	if err != nil {
		r.Logger.Error("Error query get user list", zap.Error(err))
		return nil, 0, err
	}

	return users, total, nil
}

func (r *userRepository) GetUserByID(ctx context.Context, id uint) (entity.User, error) {
	db := infra.GetDB(ctx, r.db)
	var user entity.User
	query := db.Model(&user).Where("id = ?", id).Limit(1)
	err := query.Find(&user).Error
	if err != nil {
		r.Logger.Error("Error query get user by id", zap.Error(err))
		return user, err
	}
	return user, nil
}

func (r *userRepository) UpdateUser(ctx context.Context, id uint, u *entity.User) error {
	db := infra.GetDB(ctx, r.db)
	err := db.Save(u).Error
	if err != nil {
		r.Logger.Error("Error query update user", zap.Error(err))
		return err
	}
	return nil
}

func (r *userRepository) DeleteUser(ctx context.Context, id uint) error {
	db := infra.GetDB(ctx, r.db)
	err := db.Delete(&entity.User{}, id).Error
	if err != nil {
		r.Logger.Error("Error query delete user", zap.Error(err))
		return err
	}
	return nil
}