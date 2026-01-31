package repository

import (
	"errors"
	"fmt"
	"project-POS-APP-golang-integer/internal/data/entity"
	"project-POS-APP-golang-integer/internal/dto/request"
	"project-POS-APP-golang-integer/pkg/utils"
	"strings"

	"go.uber.org/zap"
	"gorm.io/gorm"
)

type ProductRepository interface {
	Create(product *entity.Product) error
	FindByID(id uint) (*entity.Product, error)
	Update(product *entity.Product) error
	SoftDelete(id uint) error
	CheckHasOrderItems(id uint) (bool, error)
	FindByNameAndCategory(name string, categoryID uint) (*entity.Product, error)
	FindAllWithFilter(req request.GetProductsRequest) ([]entity.Product, int64, error) // 🔥 TAMBAH
	GetSoldCount(productID uint) (int64, error)                                        // 🔥 TAMBAH
}

type productRepository struct {
	db  *gorm.DB
	log *zap.Logger // 🔥 FIX: lowercase log untuk konsisten
}

func NewProductRepository(db *gorm.DB, log *zap.Logger) ProductRepository {
	return &productRepository{
		db:  db,
		log: log.With(zap.String("repository", "product")), // 🔥 FIX
	}
}

func (r *productRepository) Create(product *entity.Product) error {
	r.log.Debug("Creating product", zap.String("name", product.Name))
	return r.db.Create(product).Error
}

func (r *productRepository) FindByID(id uint) (*entity.Product, error) {
	r.log.Debug("Finding product by ID", zap.Uint("id", id))

	var product entity.Product
	err := r.db.Preload("Category").Where("id = ? AND deleted_at IS NULL", id).First(&product).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			r.log.Debug("Product not found", zap.Uint("id", id))
			return nil, nil
		}
		r.log.Error("Failed to find product", zap.Uint("id", id), zap.Error(err))
		return nil, err
	}

	r.log.Debug("Product found", zap.Uint("id", id), zap.String("name", product.Name))
	return &product, nil
}

func (r *productRepository) Update(product *entity.Product) error {
	r.log.Debug("Updating product", zap.Uint("id", product.ID), zap.String("name", product.Name))
	return r.db.Save(product).Error
}

func (r *productRepository) SoftDelete(id uint) error {
	r.log.Info("Soft deleting product", zap.Uint("id", id))

	result := r.db.Model(&entity.Product{}).Where("id = ?", id).Update("deleted_at", gorm.Expr("NOW()"))
	if result.Error != nil {
		r.log.Error("Failed to soft delete product", zap.Uint("id", id), zap.Error(result.Error))
		return result.Error
	}

	if result.RowsAffected == 0 {
		r.log.Warn("Product not found for soft delete", zap.Uint("id", id))
		return utils.ErrProductNotFound
	}

	r.log.Info("Product soft deleted successfully", zap.Uint("id", id))
	return nil
}

func (r *productRepository) CheckHasOrderItems(id uint) (bool, error) {
	r.log.Debug("Checking if product has order items", zap.Uint("id", id))

	var count int64
	err := r.db.Model(&entity.OrderItem{}).Where("product_id = ?", id).Count(&count).Error
	if err != nil {
		r.log.Error("Failed to check product order items", zap.Uint("id", id), zap.Error(err))
		return false, err
	}

	hasItems := count > 0
	r.log.Debug("Product order items check result",
		zap.Uint("id", id),
		zap.Bool("has_order_items", hasItems),
		zap.Int64("count", count))

	return hasItems, nil
}

func (r *productRepository) FindByNameAndCategory(name string, categoryID uint) (*entity.Product, error) {
	r.log.Debug("Finding product by name and category",
		zap.String("name", name),
		zap.Uint("category_id", categoryID))

	var product entity.Product
	err := r.db.Where("name = ? AND category_id = ? AND deleted_at IS NULL", name, categoryID).First(&product).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			r.log.Debug("Product not found with name and category",
				zap.String("name", name),
				zap.Uint("category_id", categoryID))
			return nil, nil
		}
		r.log.Error("Failed to find product by name and category",
			zap.String("name", name),
			zap.Uint("category_id", categoryID),
			zap.Error(err))
		return nil, err
	}

	r.log.Debug("Product found by name and category",
		zap.Uint("id", product.ID),
		zap.String("name", product.Name))
	return &product, nil
}

// FindAllWithFilter - Get products with advanced filtering and sorting
func (r *productRepository) FindAllWithFilter(req request.GetProductsRequest) ([]entity.Product, int64, error) {
	r.log.Debug("Finding products with filter",
		zap.Int("page", req.Page),
		zap.Int("limit", req.Limit),
		zap.String("status", req.Status),
		zap.String("search", req.Search),
		zap.Uint("category_id", req.CategoryID))

	var products []entity.Product
	var total int64

	// Build base query
	query := r.db.Model(&entity.Product{}).
		Preload("Category").
		Where("products.deleted_at IS NULL")

	// Apply filters
	query = r.applyProductFilters(query, req)

	// Count total before pagination
	if err := query.Count(&total).Error; err != nil {
		r.log.Error("Failed to count products", zap.Error(err))
		return nil, 0, err
	}

	// Apply sorting
	query = r.applyProductSorting(query, req)

	// Apply pagination
	offset := (req.Page - 1) * req.Limit
	err := query.Offset(offset).Limit(req.Limit).Find(&products).Error
	if err != nil {
		r.log.Error("Failed to find products", zap.Error(err))
		return nil, 0, err
	}

	r.log.Debug("Products retrieved with filter",
		zap.Int("count", len(products)),
		zap.Int64("total", total))

	return products, total, nil
}

// GetSoldCount - Get total sold quantity for a product
func (r *productRepository) GetSoldCount(productID uint) (int64, error) {
	r.log.Debug("Getting sold count for product", zap.Uint("product_id", productID))

	var totalSold int64
	err := r.db.Model(&entity.OrderItem{}).
		Select("COALESCE(SUM(quantity), 0)").
		Where("product_id = ?", productID).
		Scan(&totalSold).Error

	if err != nil {
		r.log.Error("Failed to get sold count",
			zap.Uint("product_id", productID),
			zap.Error(err))
		return 0, err
	}

	r.log.Debug("Sold count retrieved",
		zap.Uint("product_id", productID),
		zap.Int64("sold_count", totalSold))

	return totalSold, nil
}

// Helper method to apply filters
func (r *productRepository) applyProductFilters(query *gorm.DB, req request.GetProductsRequest) *gorm.DB {
	// Filter by category ID
	if req.CategoryID > 0 {
		query = query.Where("products.category_id = ?", req.CategoryID)
	}

	// Filter by category name (JOIN with categories)
	if req.CategoryName != "" {
		query = query.Joins("INNER JOIN categories ON categories.id = products.category_id").
			Where("categories.deleted_at IS NULL").
			Where("LOWER(categories.name) LIKE ?", "%"+strings.ToLower(req.CategoryName)+"%")
	}

	// Filter by status
	if req.Status != "" {
		query = query.Where("products.status = ?", req.Status)
	}

	// Filter by price range
	if req.MinPrice > 0 {
		query = query.Where("products.price >= ?", req.MinPrice)
	}
	if req.MaxPrice > 0 {
		query = query.Where("products.price <= ?", req.MaxPrice)
	}

	// Search in name or description
	if req.Search != "" {
		searchTerm := "%" + strings.ToLower(req.Search) + "%"
		query = query.Where("LOWER(products.name) LIKE ? OR LOWER(products.description) LIKE ?",
			searchTerm, searchTerm)
	}

	return query
}

// Helper method to apply sorting
func (r *productRepository) applyProductSorting(query *gorm.DB, req request.GetProductsRequest) *gorm.DB {
	// Default sorting by ID descending
	query = query.Order("products.id DESC")

	// Apply custom sorting (exclude sort_by_sold karena butuh special handling)
	if req.SortByStock != "" {
		query = query.Order(fmt.Sprintf("products.stock %s", strings.ToUpper(req.SortByStock)))
	}

	if req.SortByPrice != "" {
		query = query.Order(fmt.Sprintf("products.price %s", strings.ToUpper(req.SortByPrice)))
	}

	if req.SortByCreatedAt != "" {
		query = query.Order(fmt.Sprintf("products.created_at %s", strings.ToUpper(req.SortByCreatedAt)))
	}

	// Note: sort_by_sold akan dihandle di service layer
	// karena butuh subquery atau post-processing

	return query
}

// OPTIONAL: Method untuk mendapatkan products dengan sold count dalam satu query
func (r *productRepository) FindAllWithSoldCount(req request.GetProductsRequest) ([]struct {
	entity.Product
	SoldCount int64 `json:"sold_count"`
}, int64, error) {
	// Advanced query dengan subquery untuk sold count
	// Implement jika butuh performance lebih baik
	return nil, 0, nil
}
