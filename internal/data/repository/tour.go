package repository

import (
	"context"
	"strings"
	"travel-api/internal/data/entity"
	"travel-api/internal/dto"

	"go.uber.org/zap"
	"gorm.io/gorm"
)

type TourRepository interface{
	FindAll(ctx context.Context, f dto.TourFilterRequest) ([]entity.Tour, int64, error)
	ScheduleByID(ctx context.Context, id uint) (*entity.Tour, error)
}

type tourRepository struct {
	db *gorm.DB
	Logger *zap.Logger
}

func NewTourRepo(db *gorm.DB, log *zap.Logger) TourRepository {
	return &tourRepository{
		db: db,
		Logger: log,
	}
}

func (r *tourRepository) FindAll(ctx context.Context, f dto.TourFilterRequest) ([]entity.Tour, int64, error) {
	var tours []entity.Tour
	var totalItems int64

	query := r.db.Model(&entity.Tour{}).
		Preload("Location").
		Preload("Images", "is_main = ?", true).
		Preload("Reviews")
	
	query = query.Preload("Schedules", func(db *gorm.DB) *gorm.DB {
		return db.Where("start_date >= CURRENT_DATE").Order("start_date asc").Limit(1)
	})

	if f.Search != "" {
		searchPattern := "%" + strings.ToLower(f.Search) + "%"
		query = query.Joins("JOIN locations l ON l.id = tours.location_id").
			Where("LOWER(tours.name) LIKE ?", searchPattern)
	}

	if f.Date != "" {
		query = query.Joins("JOIN tour_schedules ts ON ts.tour_id = tours.id").
			Where("ts.start_date >= ?", f.Date).Group("tours.id")
	}

	if err := query.Count(&totalItems).Error; err != nil {
		return nil, 0, err
	}

	switch f.SortBy {
	case "price_low":
		query = query.Order("tours.base_price asc")
	case "price_hight":
		query = query.Order("tours.base_price desc")
	default:
		query = query.Order("tours.created_at desc")
	}

	offset := (f.Page - 1) * f.Limit

	err := query.Limit(f.Limit).Offset(offset).Find(&tours).Error

	return tours, totalItems, err
}

func (r *tourRepository) ScheduleByID(ctx context.Context, id uint) (*entity.Tour, error) {
	var tour entity.Tour

	query := r.db.Model(&tour).
		Preload("Location").
		Preload("Images").
		Preload("Reviews")
	
	query = query.Preload("Schedules", func(db *gorm.DB) *gorm.DB {
		return db.Where("id = ?", id).Limit(1)
	})
	
	err := query.Find(&tour).Error
	if err != nil {
		return nil, err
	}
	
	return &tour, err
}