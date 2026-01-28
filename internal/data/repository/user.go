package repository

import (
	"context"
	"project-POS-APP-golang-integer/internal/data/entity"
	"project-POS-APP-golang-integer/internal/dto"

	"go.uber.org/zap"
	"gorm.io/gorm"
)

type UserRepository interface {
	Create(ctx context.Context, user *entity.User) (*entity.User, error)
	FindByEmail(ctx context.Context, email string) (*entity.User, error)
	GetUserList(ctx context.Context, f dto.UserFilterRequest) ([]entity.User, int64, error)
	GetByID(ctx context.Context, id uint) (entity.User, error)
	Update(ctx context.Context, id uint, data *entity.User) error
	Delete(ctx context.Context, id uint) error
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

func (r *userRepository) Create(ctx context.Context, user *entity.User) (*entity.User, error) {
	tx := r.db.Begin()
	if tx.Error != nil {
		return nil, tx.Error
	}
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		} else {
			tx.Commit()
		}
	}()
	err := tx.Create(&user).Error
	if err != nil {
		r.Logger.Error("Error query create user", zap.Error(err))
		return nil, err
	}
	if user.Role == "staff" {
		var staff entity.Staff
		err = tx.Create(&staff).Error
		if err != nil {
			r.Logger.Error("Error query create staff", zap.Error(err))
			return nil, err
		}
	}
	return user, err
}

func (r *userRepository) FindByEmail(ctx context.Context, email string) (*entity.User, error) {
	var user entity.User
	query := r.db.WithContext(ctx).Model(&user).Where("email = ?", email).Limit(1)
	err := query.Find(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, err
}

func (r *userRepository) GetUserList(ctx context.Context, f dto.UserFilterRequest) ([]entity.User, int64, error) {
	var users []entity.User
	var total int64
	query := r.db.Model(&entity.User{})

	// Filter by role
	switch f.Role {
	case "admin":
		query = query.Where("role = ?", f.Role)
	case "staff":
		query = query.Preload("Staff").Where("role = ?", f.Role)
	}

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

func (r *userRepository) GetByID(ctx context.Context, id uint) (entity.User, error) {
	var user entity.User
	query := r.db.Model(&user).Where("id = ?", id).Limit(1)
	err := query.Find(&user).Error
	if err != nil {
		r.Logger.Error("Error query get user by id", zap.Error(err))
		return user, err
	}
	return user, nil
}

func (r *userRepository) Update(ctx context.Context, id uint, u *entity.User) error {
	tx := r.db.Begin()
	if tx.Error != nil {
		return tx.Error
	}
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		} else {
			tx.Commit()
		}
	}()
	err := tx.Save(u).Error
	if err != nil {
		r.Logger.Error("Error query update user", zap.Error(err))
		return err
	}
	return nil
}

func (r *userRepository) Delete(ctx context.Context, id uint) error {
	err := r.db.Delete(&entity.User{}, id).Error
	if err != nil {
		r.Logger.Error("Error query delete user", zap.Error(err))
		return err
	}
	return nil
}