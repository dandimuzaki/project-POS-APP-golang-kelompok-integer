package repository

import (
	"project-POS-APP-golang-integer/internal/data/entity"

	"gorm.io/gorm"
)

type CategoryRepository interface {
	FindAll() ([]entity.Category, error)
}

type categoryRepository struct {
	db *gorm.DB
}

func NewCategoryRepository(db *gorm.DB) CategoryRepository {
	return &categoryRepository{db}
}

func (r *categoryRepository) FindAll() ([]entity.Category, error) {
	var categories []entity.Category
	err := r.db.Find(&categories).Error
	return categories, err
}
