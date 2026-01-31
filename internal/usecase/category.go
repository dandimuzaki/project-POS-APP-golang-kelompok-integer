package usecase

import (
	"errors"
	"project-POS-APP-golang-integer/internal/data/entity"
	"project-POS-APP-golang-integer/internal/data/repository"
	"project-POS-APP-golang-integer/internal/dto/request"
	"project-POS-APP-golang-integer/internal/dto/response"
	"project-POS-APP-golang-integer/pkg/utils"

	"go.uber.org/zap"
)

type CategoryService interface {
	GetAllCategories(req request.GetCategoriesRequest) (*response.CategoryListResponse, error)
	GetCategoryByID(id uint) (*response.CategoryResponse, error)
	CreateCategory(req request.CreateCategoryRequest) (*response.CategoryResponse, error)
	UpdateCategory(id uint, req request.UpdateCategoryRequest) (*response.CategoryResponse, error)
	DeleteCategory(id uint) error
}

type categoryService struct {
	tx           TxManager
	categoryRepo repository.CategoryRepository
	log          *zap.Logger
}

func NewCategoryService(
	tx TxManager,
	categoryRepo repository.CategoryRepository,
	log *zap.Logger,
) CategoryService {
	return &categoryService{
		tx:           tx,
		categoryRepo: categoryRepo,
		log:          log.With(zap.String("service", "category")),
	}
}

func (cs *categoryService) GetAllCategories(req request.GetCategoriesRequest) (*response.CategoryListResponse, error) {
	cs.log.Debug("Getting categories",
		zap.String("name_filter", req.Name),
		zap.Int("page", req.Page),
		zap.Int("per_page", req.PerPage))

	categories, total, err := cs.categoryRepo.FindAllWithPagination(req.Name, req.Page, req.PerPage)
	if err != nil {
		cs.log.Error("Failed to get categories from repository", zap.Error(err))
		return nil, err
	}

	// Map entities to responses
	var categoryResponses []response.CategoryResponse
	for _, category := range categories {
		categoryResponses = append(categoryResponses, response.CategoryResponse{
			ID:          category.ID,
			Name:        category.Name,
			IconURL:     category.IconURL,
			Description: category.Description,
			CreatedAt:   category.CreatedAt,
			UpdatedAt:   category.UpdatedAt,
		})
	}

	// Calculate total pages
	totalPages := int(total) / req.PerPage
	if int(total)%req.PerPage > 0 {
		totalPages++
	}

	cs.log.Debug("Categories retrieved",
		zap.Int("count", len(categoryResponses)),
		zap.Int64("total", total))

	return &response.CategoryListResponse{
		Data: categoryResponses,
		Pagination: response.PaginationMeta{
			Page:       req.Page,
			PerPage:    req.PerPage,
			Total:      total,
			TotalPages: totalPages,
		},
	}, nil
}

func (cs *categoryService) GetCategoryByID(id uint) (*response.CategoryResponse, error) {
	cs.log.Debug("Getting category by ID", zap.Uint("id", id))

	category, err := cs.categoryRepo.FindByID(id)
	if err != nil {
		cs.log.Error("Failed to get category from repository", zap.Error(err))
		return nil, err
	}

	if category == nil {
		cs.log.Warn("Category not found", zap.Uint("id", id))
		return nil, utils.ErrCategoryNotFound
	}

	cs.log.Debug("Category found", zap.Uint("id", id), zap.String("name", category.Name))

	return &response.CategoryResponse{
		ID:          category.ID,
		Name:        category.Name,
		IconURL:     category.IconURL,
		Description: category.Description,
		CreatedAt:   category.CreatedAt,
		UpdatedAt:   category.UpdatedAt,
	}, nil
}

func (cs *categoryService) CreateCategory(req request.CreateCategoryRequest) (*response.CategoryResponse, error) {
	cs.log.Info("Creating new category", zap.String("name", req.Name))

	// Check if category name already exists
	existingCategory, err := cs.categoryRepo.FindByName(req.Name)
	if err != nil {
		cs.log.Error("Failed to check existing category", zap.Error(err))
		return nil, err
	}

	if existingCategory != nil {
		cs.log.Warn("Category name already exists", zap.String("name", req.Name))
		return nil, utils.ErrCategoryExists
	}

	// Create new category entity
	category := &entity.Category{
		Name:        req.Name,
		IconURL:     req.IconURL,
		Description: req.Description,
	}

	err = cs.categoryRepo.Create(category)
	if err != nil {
		cs.log.Error("Failed to create category in repository", zap.Error(err))
		return nil, err
	}

	cs.log.Info("Category created successfully",
		zap.Uint("id", category.ID),
		zap.String("name", category.Name))

	return &response.CategoryResponse{
		ID:          category.ID,
		Name:        category.Name,
		IconURL:     category.IconURL,
		Description: category.Description,
		CreatedAt:   category.CreatedAt,
		UpdatedAt:   category.UpdatedAt,
	}, nil
}

func (cs *categoryService) UpdateCategory(id uint, req request.UpdateCategoryRequest) (*response.CategoryResponse, error) {
	cs.log.Info("Updating category", zap.Uint("id", id))

	// Get existing category
	category, err := cs.categoryRepo.FindByID(id)
	if err != nil {
		cs.log.Error("Failed to get category for update", zap.Error(err))
		return nil, err
	}

	if category == nil {
		cs.log.Warn("Category not found for update", zap.Uint("id", id))
		return nil, utils.ErrCategoryNotFound
	}

	// Check if new name already exists (if name is being changed)
	if req.Name != "" && req.Name != category.Name {
		existingCategory, err := cs.categoryRepo.FindByName(req.Name)
		if err != nil {
			cs.log.Error("Failed to check existing category name", zap.Error(err))
			return nil, err
		}
		if existingCategory != nil {
			cs.log.Warn("Category name already exists", zap.String("name", req.Name))
			return nil, utils.ErrCategoryExists
		}
		category.Name = req.Name
	}
	if req.Name == "" && req.IconURL == "" && req.Description == "" {
		cs.log.Warn("No changes provided for update", zap.Uint("id", id))
		return nil, errors.New("no changes provided")
	}

	// Update other fields if provided
	if req.IconURL != "" {
		category.IconURL = req.IconURL
	}
	if req.Description != "" {
		category.Description = req.Description
	}

	err = cs.categoryRepo.Update(category)
	if err != nil {
		cs.log.Error("Failed to update category in repository", zap.Error(err))
		return nil, err
	}

	cs.log.Info("Category updated successfully", zap.Uint("id", id))

	return &response.CategoryResponse{
		ID:          category.ID,
		Name:        category.Name,
		IconURL:     category.IconURL,
		Description: category.Description,
		CreatedAt:   category.CreatedAt,
		UpdatedAt:   category.UpdatedAt,
	}, nil
}

func (cs *categoryService) DeleteCategory(id uint) error {
	cs.log.Info("Deleting category", zap.Uint("id", id))

	// Check if category exists
	category, err := cs.categoryRepo.FindByID(id)
	if err != nil {
		cs.log.Error("Failed to get category for deletion", zap.Error(err))
		return err
	}

	if category == nil {
		cs.log.Warn("Category not found for deletion", zap.Uint("id", id))
		return utils.ErrCategoryNotFound
	}

	// Check if category has associated products
	hasProducts, err := cs.categoryRepo.CheckHasProducts(id)
	if err != nil {
		cs.log.Error("Failed to check category products", zap.Error(err))
		return err
	}

	if hasProducts {
		cs.log.Warn("Cannot delete category with associated products", zap.Uint("id", id))
		return utils.ErrCategoryHasProducts
	}

	// Perform soft delete
	err = cs.categoryRepo.SoftDelete(id)
	if err != nil {
		cs.log.Error("Failed to delete category from repository", zap.Error(err))
		return err
	}

	cs.log.Info("Category deleted successfully", zap.Uint("id", id))
	return nil
}
