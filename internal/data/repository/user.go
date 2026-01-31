package repository

import (
	"context"
	"errors"
	"project-POS-APP-golang-integer/internal/data/entity"
	"project-POS-APP-golang-integer/internal/infra"
	"project-POS-APP-golang-integer/pkg/utils"
	"strings"

	"go.uber.org/zap"
	"gorm.io/gorm"
)

type UserRepository interface {
	CreateUser(ctx context.Context, user *entity.User) (*entity.User, error)
	FindUserByEmail(ctx context.Context, email string) (*entity.User, error)
	GetUserList(ctx context.Context, f UserQueryParams) ([]entity.User, int64, error)
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

type UserQueryParams struct {
	Offset int
	Limit int
	Role entity.UserRole
	Name string
	Email string
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

func (r *userRepository) GetUserList(ctx context.Context, f UserQueryParams) ([]entity.User, int64, error) {
	db := infra.GetDB(ctx, r.db)
	var users []entity.User
	var total int64
	query := db.Model(&entity.User{})

	// Filter by role
	if f.Role != "" {
		query = query.Where("role = ?", f.Role)
	}
	
	// Search by name
	if f.Name != "" {
		searchPattern := "%" + strings.ToLower(f.Name) + "%"
		query = query.Joins("JOIN profiles p ON p.user_id = users.id").
			Where("LOWER(full_name) LIKE ?", searchPattern)
	}

	// Search by email
	if f.Email != "" {
		searchPattern := "%" + strings.ToLower(f.Email) + "%"
		query = query.Where("LOWER(email) LIKE ?", searchPattern)
	}

	// Get total user
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	err := query.Limit(f.Limit).Offset(f.Offset).Find(&users).Error
	if err != nil {
		r.Logger.Error("Error query get user list", zap.Error(err))
		return nil, 0, err
	}

	return users, total, nil
}

func (r *userRepository) GetUserByID(ctx context.Context, id uint) (entity.User, error) {
	db := infra.GetDB(ctx, r.db)

	var user entity.User
	err := db.
		Model(&user).
		Where("id = ?", id).
		Limit(1).
		First(&user).
		Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return user, utils.ErrUserNotFound
		}

		r.Logger.Error("Error query get user by id", zap.Error(err))
		return user, err
	}

	return user, nil
}

func (r *userRepository) UpdateUser(ctx context.Context, id uint, u *entity.User) error {
	db := infra.GetDB(ctx, r.db)

	result := db.Model(&entity.User{}).
		Where("id = ?", u.ID).
		Updates(u)

	if result.Error != nil {
		r.Logger.Error("Error query update user", zap.Error(result.Error))
		return result.Error
	}

	if result.RowsAffected == 0 {
		return utils.ErrUserNotFound
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