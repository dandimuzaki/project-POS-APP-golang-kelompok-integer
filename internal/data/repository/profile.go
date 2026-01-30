package repository

import (
	"context"
	"project-POS-APP-golang-integer/internal/data/entity"
	"project-POS-APP-golang-integer/internal/infra"

	"go.uber.org/zap"
	"gorm.io/gorm"
)

type ProfileRepository interface {
	GetProfileByID(ctx context.Context, id uint) (*entity.User, error)
	CreateProfile(ctx context.Context, profile *entity.Profile) (*entity.Profile, error)
	UpdateProfile(ctx context.Context, data *entity.Profile) error
}

type profileRepository struct {
	db *gorm.DB
	Logger *zap.Logger
}

func NewProfileRepo(db *gorm.DB, log *zap.Logger) ProfileRepository {
	return &profileRepository{
		db: db,
		Logger: log,
	}
}

func (r *profileRepository) GetProfileByID(ctx context.Context, id uint) (*entity.User, error) {
	db := infra.GetDB(ctx, r.db)
	var user entity.User
	query := db.Model(&user).Where("id = ?", id).Limit(1).Preload("Profile")
	err := query.Find(&user).Error
	if err != nil {
		r.Logger.Error("Error query get profile by id", zap.Error(err))
		return &user, err
	}
	return &user, nil
}

func (r *profileRepository) CreateProfile(ctx context.Context, profile *entity.Profile) (*entity.Profile, error) {
	db := infra.GetDB(ctx, r.db)
	err := db.Create(&profile).Error
	if err != nil {
		r.Logger.Error("Error query create profile", zap.Error(err))
		return nil, err
	}
	return profile, err
}

func (r *profileRepository) UpdateProfile(ctx context.Context, p *entity.Profile) error {
	db := infra.GetDB(ctx, r.db)
	err := db.Model(&entity.Profile{}).
    Where("user_id = ?", p.UserID).
    Updates(p).Error
	if err != nil {
		r.Logger.Error("Error query update profile", zap.Error(err))
		return err
	}
	return nil
}