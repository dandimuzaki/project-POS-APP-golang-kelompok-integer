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

type ProductHandler struct {
	srv usecase.ProductService
	log *zap.Logger
}

func NewProductHandler(
	srv usecase.ProductService,
	log *zap.Logger,
) *ProductHandler {
	return &ProductHandler{
		srv: srv,
		log: log.With(zap.String("handler", "product")),
	}
}

// CreateProduct creates a new product
// @Summary Create a new product
// @Description Create a new product
// @Tags Products
// @Accept json
// @Produce json
// @Param request body request.CreateProductRequest true "Product data"
// @Success 201 {object} response.SuccessResponse
// @Failure 400 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /api/v1/products [post]
func (h *ProductHandler) CreateProduct(c *gin.Context) {
	var req request.CreateProductRequest

	// Bind JSON
	if err := c.ShouldBindJSON(&req); err != nil {
		h.log.Warn("Failed to bind request body", zap.Error(err))
		utils.ResponseFailed(c, 400, "Invalid request body", err.Error())
		return
	}

	// Validate
	if err := req.Validate(); err != nil {
		h.log.Warn("Request validation failed", zap.Error(err))
		utils.ResponseFailed(c, 400, "Validation failed", err.Error())
		return
	}

	// Call service
	product, err := h.srv.CreateProduct(req)
	if err != nil {
		h.log.Error("Failed to create product", zap.Error(err))

		if utils.IsBusinessError(err) {
			utils.ResponseFailed(c, 400, err.Error(), nil)
		} else {
			utils.ResponseFailed(c, 500, "Failed to create product", nil)
		}
		return
	}

	h.log.Info("Product created successfully",
		zap.Uint("product_id", product.ID),
		zap.String("name", product.Name))

	// Created Response (201)
	utils.ResponseSuccess(c, 201, "Product created successfully", product)
}

// GetAllProducts gets list of products with filters
// @Summary Get all products
// @Description Get list of all products with advanced filtering and sorting
// @Tags Products
// @Accept json
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Items per page" default(10)
// @Param category_id query int false "Filter by category ID"
// @Param category_name query string false "Filter by category name"
// @Param status query string false "Filter by status (active/inactive)"
// @Param min_price query number false "Minimum price"
// @Param max_price query number false "Maximum price"
// @Param search query string false "Search in product name or description"
// @Param sort_by_stock query string false "Sort by stock (asc/desc)"
// @Param sort_by_price query string false "Sort by price (asc/desc)"
// @Param sort_by_created_at query string false "Sort by created date (asc/desc)"
// @Param sort_by_sold query string false "Sort by sold count (asc/desc)"
// @Success 200 {object} response.SuccessResponse
// @Failure 400 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /api/v1/products [get]
func (h *ProductHandler) GetAllProducts(c *gin.Context) {
	var req request.GetProductsRequest

	// Bind query parameters
	if err := c.ShouldBindQuery(&req); err != nil {
		h.log.Warn("Failed to bind query parameters", zap.Error(err))
		utils.ResponseFailed(c, 400, "Invalid query parameters", err.Error())
		return
	}

	// Set default values
	if req.Page == 0 {
		req.Page = 1
	}
	if req.Limit == 0 {
		req.Limit = 10
	}

	// Validate pagination
	if req.Page <= 0 {
		h.log.Warn("Invalid page number", zap.Int("page", req.Page))
		utils.ResponseFailed(c, 400, "Page must be greater than 0", nil)
		return
	}

	if req.Limit <= 0 || req.Limit > 100 {
		h.log.Warn("Invalid limit value", zap.Int("limit", req.Limit))
		utils.ResponseFailed(c, 400, "Limit must be between 1 and 100", nil)
		return
	}

	// Validate sort parameters
	if req.SortByStock != "" && req.SortByStock != "asc" && req.SortByStock != "desc" {
		h.log.Warn("Invalid sort_by_stock value", zap.String("sort_by_stock", req.SortByStock))
		utils.ResponseFailed(c, 400, "sort_by_stock must be 'asc' or 'desc'", nil)
		return
	}

	if req.SortByPrice != "" && req.SortByPrice != "asc" && req.SortByPrice != "desc" {
		h.log.Warn("Invalid sort_by_price value", zap.String("sort_by_price", req.SortByPrice))
		utils.ResponseFailed(c, 400, "sort_by_price must be 'asc' or 'desc'", nil)
		return
	}

	if req.SortByCreatedAt != "" && req.SortByCreatedAt != "asc" && req.SortByCreatedAt != "desc" {
		h.log.Warn("Invalid sort_by_created_at value", zap.String("sort_by_created_at", req.SortByCreatedAt))
		utils.ResponseFailed(c, 400, "sort_by_created_at must be 'asc' or 'desc'", nil)
		return
	}

	if req.SortBySold != "" && req.SortBySold != "asc" && req.SortBySold != "desc" {
		h.log.Warn("Invalid sort_by_sold value", zap.String("sort_by_sold", req.SortBySold))
		utils.ResponseFailed(c, 400, "sort_by_sold must be 'asc' or 'desc'", nil)
		return
	}

	// Validate price range
	if req.MinPrice > 0 && req.MaxPrice > 0 && req.MinPrice > req.MaxPrice {
		h.log.Warn("Invalid price range",
			zap.Float64("min_price", req.MinPrice),
			zap.Float64("max_price", req.MaxPrice))
		utils.ResponseFailed(c, 400, "min_price cannot be greater than max_price", nil)
		return
	}

	h.log.Debug("GetAllProducts request",
		zap.Int("page", req.Page),
		zap.Int("limit", req.Limit),
		zap.String("search", req.Search),
		zap.String("status", req.Status),
		zap.String("client_ip", c.ClientIP()))

	// Call service
	result, err := h.srv.GetAllProducts(req)
	if err != nil {
		h.log.Error("Failed to get products", zap.Error(err))
		utils.ResponseFailed(c, 500, "Failed to retrieve products", nil)
		return
	}

	h.log.Debug("Products retrieved successfully",
		zap.Int("count", len(result.Data)),
		zap.Int64("total", result.Pagination.Total))

	// Success Response (200 OK)
	utils.ResponsePagination(c, 200, "Products retrieved successfully",
		result.Data, response.PaginationMeta{
			Page:       result.Pagination.Page,
			PerPage:    result.Pagination.PerPage,
			Total:      result.Pagination.Total,
			TotalPages: result.Pagination.TotalPages,
		})
}

// GetProductByID gets a product by ID
// @Summary Get product by ID
// @Description Get product details by ID
// @Tags Products
// @Accept json
// @Produce json
// @Param id path string true "Product ID"
// @Success 200 {object} response.SuccessResponse
// @Failure 400 {object} response.ErrorResponse
// @Failure 404 {object} response.ErrorResponse
// @Router /api/v1/products/{id} [get]
func (h *ProductHandler) GetProductByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil || id == 0 {
		h.log.Warn("Invalid product ID",
			zap.String("id", idStr),
			zap.Error(err))
		utils.ResponseFailed(c, 400, "Invalid product ID", nil)
		return
	}

	// Call service
	product, err := h.srv.GetProductByID(uint(id))
	if err != nil {
		h.log.Error("Failed to get product", zap.Uint("id", uint(id)), zap.Error(err))

		if err == utils.ErrProductNotFound {
			utils.ResponseFailed(c, 404, "Product not found", nil)
		} else {
			utils.ResponseFailed(c, 500, "Failed to get product", nil)
		}
		return
	}

	// Success Response (200 OK)
	utils.ResponseSuccess(c, 200, "Product retrieved successfully", product)
}

// UpdateProduct updates a product
// @Summary Update a product
// @Description Update product details
// @Tags Products
// @Accept json
// @Produce json
// @Param id path string true "Product ID"
// @Param request body request.UpdateProductRequest true "Product data"
// @Success 200 {object} response.SuccessResponse
// @Failure 400 {object} response.ErrorResponse
// @Failure 404 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /api/v1/products/{id} [put]
func (h *ProductHandler) UpdateProduct(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil || id == 0 {
		h.log.Warn("Invalid product ID", zap.String("id", idStr), zap.Error(err))
		utils.ResponseFailed(c, 400, "Invalid product ID", nil)
		return
	}

	var req request.UpdateProductRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.log.Warn("Failed to bind request body", zap.Error(err))
		utils.ResponseFailed(c, 400, "Invalid request body", err.Error())
		return
	}

	// Validate
	if err := req.Validate(); err != nil {
		h.log.Warn("Request validation failed", zap.Error(err))
		utils.ResponseFailed(c, 400, "Validation failed", err.Error())
		return
	}

	// Check if no changes provided
	if req.Name == "" && req.Description == "" && req.Price == 0 &&
		req.Stock == 0 && req.MinStock == 0 && req.ImageURL == "" &&
		req.Status == "" && req.CategoryID == 0 {
		h.log.Warn("No changes provided for update", zap.Uint("id", uint(id)))
		utils.ResponseFailed(c, 400, "No changes provided", nil)
		return
	}

	// Call service
	product, err := h.srv.UpdateProduct(uint(id), req)
	if err != nil {
		h.log.Error("Failed to update product",
			zap.Uint("id", uint(id)),
			zap.Error(err))

		if err == utils.ErrProductNotFound {
			utils.ResponseFailed(c, 404, "Product not found", nil)
		} else if err == utils.ErrCategoryNotFound {
			utils.ResponseFailed(c, 400, "Category not found", nil)
		} else if err == utils.ErrNoChangesProvided {
			utils.ResponseFailed(c, 400, "No changes provided", nil)
		} else if utils.IsBusinessError(err) {
			utils.ResponseFailed(c, 400, err.Error(), nil)
		} else {
			utils.ResponseFailed(c, 500, "Failed to update product", nil)
		}
		return
	}

	h.log.Info("Product updated successfully",
		zap.Uint("id", uint(id)))

	// Success Response (200 OK)
	utils.ResponseSuccess(c, 200, "Product updated successfully", product)
}

// DeleteProduct deletes a product (soft delete)
// @Summary Delete a product
// @Description Soft delete a product
// @Tags Products
// @Accept json
// @Produce json
// @Param id path string true "Product ID"
// @Success 200 {object} response.SuccessResponse
// @Failure 400 {object} response.ErrorResponse
// @Failure 404 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /api/v1/products/{id} [delete]
func (h *ProductHandler) DeleteProduct(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil || id == 0 {
		h.log.Warn("Invalid product ID", zap.String("id", idStr), zap.Error(err))
		utils.ResponseFailed(c, 400, "Invalid product ID", nil)
		return
	}

	// Call service
	if err := h.srv.DeleteProduct(uint(id)); err != nil {
		h.log.Error("Failed to delete product", zap.Uint("id", uint(id)), zap.Error(err))

		if err == utils.ErrProductNotFound {
			utils.ResponseFailed(c, 404, "Product not found", nil)
		} else if err == utils.ErrProductHasOrders {
			utils.ResponseFailed(c, 400, "Cannot delete product with associated orders", nil)
		} else if utils.IsBusinessError(err) {
			utils.ResponseFailed(c, 400, err.Error(), nil)
		} else {
			utils.ResponseFailed(c, 500, "Failed to delete product", nil)
		}
		return
	}

	h.log.Info("Product deleted successfully", zap.Uint("id", uint(id)))

	// Success Response (200 OK)
	utils.ResponseSuccess(c, 200, "Product deleted successfully", nil)
}
