package adaptor

import (
	"project-POS-APP-golang-integer/internal/dto/request"
	"project-POS-APP-golang-integer/internal/dto/response"
	"project-POS-APP-golang-integer/internal/usecase"
	"project-POS-APP-golang-integer/pkg/utils"
	"strconv"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type CategoryHandler struct {
	srv usecase.CategoryService
	log *zap.Logger
}

func NewCategoryHandler(
	srv usecase.CategoryService,
	log *zap.Logger,
) *CategoryHandler {
	return &CategoryHandler{
		srv: srv,
		log: log.With(zap.String("handler", "category")),
	}
}

// GetAllCategories gets list of categories with filters
// @Summary Get all categories
// @Description Get list of all categories with optional filters
// @Tags Categories
// @Accept json
// @Produce json
// @Param name query string false "Filter by category name"
// @Param page query int false "Page number" default(1)
// @Param per_page query int false "Items per page" default(10)
// @Success 200 {object} response.SuccessResponse
// @Failure 400 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /api/v1/categories [get]
func (ch *CategoryHandler) GetAllCategories(c *gin.Context) {
	var req request.GetCategoriesRequest

	// Bind query parameters
	if err := c.ShouldBindQuery(&req); err != nil {
		ch.log.Warn("Failed to bind query parameters", zap.Error(err))
		utils.ResponseFailed(c, 400, "Invalid query parameters", err.Error())
		return
	}

	// Set default values
	if req.Page == 0 {
		req.Page = 1
	}
	if req.PerPage == 0 {
		req.PerPage = 10
	}

	// Validate pagination (Bad Request example)
	if req.Page <= 0 {
		ch.log.Warn("Invalid page number", zap.Int("page", req.Page))
		utils.ResponseFailed(c, 400, "Page must be greater than 0", nil)
		return
	}

	if req.PerPage <= 0 || req.PerPage > 100 {
		ch.log.Warn("Invalid per_page value", zap.Int("per_page", req.PerPage))
		utils.ResponseFailed(c, 400, "Per page must be between 1 and 100", nil)
		return
	}

	// Call service
	result, err := ch.srv.GetAllCategories(req)
	if err != nil {
		ch.log.Error("Failed to get categories", zap.Error(err))
		utils.ResponseFailed(c, 500, "Failed to retrieve categories", nil)
		return
	}

	ch.log.Debug("Categories retrieved successfully",
		zap.Int("count", len(result.Data)),
		zap.Int64("total", result.Pagination.Total))

	// Success Response (200 OK)
	utils.ResponsePagination(c, 200, "Categories retrieved successfully",
		result.Data, response.PaginationMeta{
			Page:       result.Pagination.Page,
			PerPage:    result.Pagination.PerPage,
			Total:      result.Pagination.Total,
			TotalPages: result.Pagination.TotalPages,
		})
}

// GetCategoryByID gets a category by ID
// @Summary Get category by ID
// @Description Get category details by ID
// @Tags Categories
// @Accept json
// @Produce json
// @Param id path string true "Category ID"
// @Success 200 {object} response.SuccessResponse
// @Failure 400 {object} response.ErrorResponse
// @Failure 404 {object} response.ErrorResponse
// @Router /api/v1/categories/{id} [get]
func (ch *CategoryHandler) GetCategoryByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		ch.log.Warn("Invalid category ID", zap.String("id", idStr), zap.Error(err))
		utils.ResponseFailed(c, 400, "Invalid category ID", nil)
		return
	}

	// Call service
	category, err := ch.srv.GetCategoryByID(uint(id))
	if err != nil {
		ch.log.Error("Failed to get category", zap.Uint("id", uint(id)), zap.Error(err))

		if err == utils.ErrCategoryNotFound {
			utils.ResponseFailed(c, 404, "Category not found", nil)
		} else {
			utils.ResponseFailed(c, 500, "Failed to get category", nil)
		}
		return
	}

	// Success Response (200 OK)
	utils.ResponseSuccess(c, 200, "Category retrieved successfully", category)
}

// CreateCategory creates a new category
// @Summary Create a new category
// @Description Create a new category
// @Tags Categories
// @Accept json
// @Produce json
// @Param request body request.CreateCategoryRequest true "Category data"
// @Success 201 {object} response.SuccessResponse
// @Failure 400 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /api/v1/categories [post]
func (ch *CategoryHandler) CreateCategory(c *gin.Context) {
	var req request.CreateCategoryRequest

	// Bind JSON
	if err := c.ShouldBindJSON(&req); err != nil {
		ch.log.Warn("Failed to bind request body", zap.Error(err))
		utils.ResponseFailed(c, 400, "Invalid request body", err.Error())
		return
	}

	// Validate
	if validationErrors, err := utils.ValidateErrors(req); err != nil {
		ch.log.Warn("Request validation failed",
			zap.Any("errors", validationErrors),
			zap.Error(err))

		utils.ResponseFailed(c, 400, "Validation failed", validationErrors)
		return
	}

	// Call service
	category, err := ch.srv.CreateCategory(req)
	if err != nil {
		ch.log.Error("Failed to create category", zap.Error(err))

		if utils.IsBusinessError(err) {
			utils.ResponseFailed(c, 400, err.Error(), nil)
		} else {
			utils.ResponseFailed(c, 500, "Failed to create category", nil)
		}
		return
	}

	ch.log.Info("Category created successfully",
		zap.Uint("category_id", category.ID))

	// Created Response (201)
	utils.ResponseSuccess(c, 201, "Category created successfully", category)
}

// UpdateCategory updates a category
// @Summary Update a category
// @Description Update category details
// @Tags Categories
// @Accept json
// @Produce json
// @Param id path string true "Category ID"
// @Param request body request.UpdateCategoryRequest true "Category data"
// @Success 200 {object} response.SuccessResponse
// @Failure 400 {object} response.ErrorResponse
// @Failure 404 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /api/v1/categories/{id} [put]
func (ch *CategoryHandler) UpdateCategory(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		ch.log.Warn("Invalid category ID", zap.String("id", idStr), zap.Error(err))
		utils.ResponseFailed(c, 400, "Invalid category ID", nil)
		return
	}

	var req request.UpdateCategoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		ch.log.Warn("Failed to bind request body", zap.Error(err))
		utils.ResponseFailed(c, 400, "Invalid request body", err.Error())
		return
	}

	// Validate
	if validationErrors, err := utils.ValidateErrors(req); err != nil {
		ch.log.Warn("Request validation failed", zap.Any("errors", validationErrors))
		utils.ResponseFailed(c, 400, "Validation failed", validationErrors)
		return
	}

	// Call service
	category, err := ch.srv.UpdateCategory(uint(id), req)
	if err != nil {
		ch.log.Error("Failed to update category",
			zap.Uint("id", uint(id)),
			zap.Error(err))

		if err == utils.ErrCategoryNotFound {
			utils.ResponseFailed(c, 404, "Category not found", nil)
		} else if utils.IsBusinessError(err) {
			utils.ResponseFailed(c, 400, err.Error(), nil)
		} else {
			utils.ResponseFailed(c, 500, "Failed to update category", nil)
		}
		return
	}

	ch.log.Info("Category updated successfully",
		zap.Uint("id", uint(id)))

	// Success Response (200 OK)
	utils.ResponseSuccess(c, 200, "Category updated successfully", category)
}

// DeleteCategory deletes a category (soft delete)
// @Summary Delete a category
// @Description Soft delete a category
// @Tags Categories
// @Accept json
// @Produce json
// @Param id path string true "Category ID"
// @Success 200 {object} response.SuccessResponse
// @Failure 400 {object} response.ErrorResponse
// @Failure 404 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /api/v1/categories/{id} [delete]
func (ch *CategoryHandler) DeleteCategory(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		ch.log.Warn("Invalid category ID", zap.String("id", idStr), zap.Error(err))
		utils.ResponseFailed(c, 400, "Invalid category ID", nil)
		return
	}

	// Call service
	if err := ch.srv.DeleteCategory(uint(id)); err != nil {
		ch.log.Error("Failed to delete category", zap.Uint("id", uint(id)), zap.Error(err))

		if err == utils.ErrCategoryNotFound {
			utils.ResponseFailed(c, 404, "Category not found", nil)
		} else if utils.IsBusinessError(err) {
			utils.ResponseFailed(c, 400, err.Error(), nil)
		} else {
			utils.ResponseFailed(c, 500, "Failed to delete category", nil)
		}
		return
	}

	ch.log.Info("Category deleted successfully", zap.Uint("id", uint(id)))

	// Success Response (200 OK)
	utils.ResponseSuccess(c, 200, "Category deleted successfully", nil)
}
