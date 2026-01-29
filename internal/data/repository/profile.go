package repository

import (
	"context"
	"project-POS-APP-golang-integer/internal/data/entity"
	"project-POS-APP-golang-integer/internal/infra"

	"go.uber.org/zap"
	"gorm.io/gorm"
)

type ProfileRepository interface {
	CreateProfile(ctx context.Context, profile *entity.Profile) (*entity.Profile, error)
	UpdateProfile(ctx context.Context, id uint, data *entity.Profile) error
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

func (r *profileRepository) CreateProfile(ctx context.Context, profile *entity.Profile) (*entity.Profile, error) {
	db := infra.GetDB(ctx, r.db)
	err := db.Create(&profile).Error
	if err != nil {
		r.Logger.Error("Error query create profile", zap.Error(err))
		return nil, err
	}
	return profile, err
}

func (r *profileRepository) UpdateProfile(ctx context.Context, id uint, u *entity.Profile) error {
	db := infra.GetDB(ctx, r.db)
	err := db.Save(u).Error
	if err != nil {
		r.Logger.Error("Error query update profile", zap.Error(err))
		return err
	}
	return nil
}