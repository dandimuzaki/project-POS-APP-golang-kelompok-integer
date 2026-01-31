package repository

import (
	"errors"
	"project-POS-APP-golang-integer/internal/data/entity"
	"project-POS-APP-golang-integer/pkg/utils"

	"go.uber.org/zap"
	"gorm.io/gorm"
)

type CategoryRepository interface {
	FindAllWithPagination(name string, page, perPage int) ([]entity.Category, int64, error)
	FindByID(id uint) (*entity.Category, error)
	FindByName(name string) (*entity.Category, error)
	Create(category *entity.Category) error
	Update(category *entity.Category) error
	SoftDelete(id uint) error
	CheckHasProducts(id uint) (bool, error)
}

type categoryRepository struct {
	db     *gorm.DB
	Logger *zap.Logger
}

func NewCategoryRepository(db *gorm.DB, log *zap.Logger) CategoryRepository {
	return &categoryRepository{
		db:     db,
		Logger: log,
	}
}

func (r *categoryRepository) FindAllWithPagination(name string, page, perPage int) ([]entity.Category, int64, error) {
	var categories []entity.Category
	var total int64

	query := r.db.Model(&entity.Category{}).Where("deleted_at IS NULL")

	if name != "" {
		query = query.Where("name LIKE ?", "%"+name+"%")
	}

	// Hitung total
	err := query.Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	// Apply pagination
	offset := (page - 1) * perPage
	err = query.Order("name ASC").Offset(offset).Limit(perPage).Find(&categories).Error

	return categories, total, err
}

func (r *categoryRepository) FindByID(id uint) (*entity.Category, error) {
	var category entity.Category
	err := r.db.Where("id = ? AND deleted_at IS NULL", id).First(&category).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil // Return nil, bukan error
		}
		return nil, err
	}
	return &category, nil
}

func (r *categoryRepository) FindByName(name string) (*entity.Category, error) {
	var category entity.Category
	err := r.db.Where("name = ? AND deleted_at IS NULL", name).First(&category).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &category, nil
}

func (r *categoryRepository) Create(category *entity.Category) error {
	return r.db.Create(category).Error
}

func (r *categoryRepository) Update(category *entity.Category) error {
	return r.db.Save(category).Error
}

func (r *categoryRepository) SoftDelete(id uint) error {
	result := r.db.Model(&entity.Category{}).Where("id = ?", id).Update("deleted_at", gorm.Expr("NOW()"))
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return utils.ErrCategoryNotFound
	}
	return nil
}

func (r *categoryRepository) CheckHasProducts(id uint) (bool, error) {
	var count int64
	err := r.db.Model(&entity.Product{}).Where("category_id = ?", id).Count(&count).Error
	return count > 0, err
}
