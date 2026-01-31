package usecase

import (
	"project-POS-APP-golang-integer/internal/data/entity"
	"project-POS-APP-golang-integer/internal/data/repository"
	"project-POS-APP-golang-integer/internal/dto/request"
	"project-POS-APP-golang-integer/internal/dto/response"
	"project-POS-APP-golang-integer/pkg/utils"

	"go.uber.org/zap"
)

type ProductService interface {
	CreateProduct(req request.CreateProductRequest) (*response.ProductResponse, error)
	GetAllProducts(req request.GetProductsRequest) (*response.ProductListResponse, error) // INI!
	GetProductByID(id uint) (*response.ProductResponse, error)
	UpdateProduct(id uint, req request.UpdateProductRequest) (*response.ProductResponse, error)
	DeleteProduct(id uint) error
}

func (s *productService) GetAllProducts(req request.GetProductsRequest) (*response.ProductListResponse, error) {
	s.log.Info("Getting products with filters",
		zap.Int("page", req.Page),
		zap.Int("limit", req.Limit))

	// Get products with filters dari repository
	products, total, err := s.productRepo.FindAllWithFilter(req)
	if err != nil {
		s.log.Error("Failed to get products", zap.Error(err))
		return nil, err
	}

	// Convert ke response
	productResponses := make([]response.ProductResponse, len(products))
	for i, product := range products {
		productResponses[i] = *response.ToProductResponse(&product)

		// // Optional: Get sold count
		// soldCount, _ := s.productRepo.GetSoldCount(product.ID)
		// productResponses[i].SoldCount = soldCount
	}

	// Calculate pagination
	totalPages := 0
	if total > 0 {
		totalPages = int(total) / req.Limit
		if int(total)%req.Limit > 0 {
			totalPages++
		}
	}

	return &response.ProductListResponse{
		Data: productResponses,
		Pagination: response.PaginationMeta{
			Page:       req.Page,
			PerPage:    req.Limit,
			Total:      total,
			TotalPages: totalPages,
		},
	}, nil
}

type productService struct {
	tx           TxManager
	productRepo  repository.ProductRepository
	categoryRepo repository.CategoryRepository
	log          *zap.Logger
}

func NewProductService(
	tx TxManager,
	productRepo repository.ProductRepository,
	categoryRepo repository.CategoryRepository,
	log *zap.Logger,
) ProductService {
	return &productService{
		tx:           tx,
		productRepo:  productRepo,
		categoryRepo: categoryRepo,
		log:          log.With(zap.String("service", "product")),
	}
}

func (s *productService) CreateProduct(req request.CreateProductRequest) (*response.ProductResponse, error) {
	s.log.Info("Creating new product", zap.String("name", req.Name))

	// Validate category exists
	category, err := s.categoryRepo.FindByID(req.CategoryID)
	if err != nil {
		s.log.Error("Failed to find category", zap.Uint("category_id", req.CategoryID), zap.Error(err))
		return nil, err
	}
	if category == nil {
		s.log.Warn("Category not found", zap.Uint("category_id", req.CategoryID))
		return nil, utils.ErrCategoryNotFound
	}

	// Check for duplicate product name in same category
	existingProduct, err := s.productRepo.FindByNameAndCategory(req.Name, req.CategoryID)
	if err != nil {
		s.log.Error("Failed to check existing product", zap.Error(err))
		return nil, err
	}
	if existingProduct != nil {
		s.log.Warn("Product name already exists in category",
			zap.String("name", req.Name),
			zap.Uint("category_id", req.CategoryID))
		return nil, utils.ErrProductExists
	}

	// Set default status if not provided
	status := entity.ProductStatusActive
	if req.Status != "" {
		status = entity.ProductStatus(req.Status)
	}

	// Create product entity
	product := &entity.Product{
		Name:        req.Name,
		Description: req.Description,
		Price:       req.Price,
		Stock:       req.Stock,
		MinStock:    req.MinStock,
		Status:      status,
		ImageURL:    req.ImageURL,
		CategoryID:  req.CategoryID,
	}

	err = s.productRepo.Create(product)
	if err != nil {
		s.log.Error("Failed to create product in repository", zap.Error(err))
		return nil, err
	}

	// Preload category for response
	product.Category = *category

	s.log.Info("Product created successfully",
		zap.Uint("product_id", product.ID),
		zap.String("name", product.Name))

	return response.ToProductResponse(product), nil
}

func (s *productService) GetProductByID(id uint) (*response.ProductResponse, error) {
	s.log.Debug("Getting product by ID", zap.Uint("id", id))

	product, err := s.productRepo.FindByID(id)
	if err != nil {
		s.log.Error("Failed to get product from repository", zap.Error(err))
		return nil, err
	}

	if product == nil {
		s.log.Warn("Product not found", zap.Uint("id", id))
		return nil, utils.ErrProductNotFound
	}

	s.log.Debug("Product found", zap.Uint("id", id), zap.String("name", product.Name))
	return response.ToProductResponse(product), nil
}

func (s *productService) UpdateProduct(id uint, req request.UpdateProductRequest) (*response.ProductResponse, error) {
	s.log.Info("Updating product", zap.Uint("id", id))

	// Get existing product
	product, err := s.productRepo.FindByID(id)
	if err != nil {
		s.log.Error("Failed to get product for update", zap.Error(err))
		return nil, err
	}

	if product == nil {
		s.log.Warn("Product not found for update", zap.Uint("id", id))
		return nil, utils.ErrProductNotFound
	}

	// Check if updating category
	if req.CategoryID != 0 && req.CategoryID != product.CategoryID {
		// Validate new category exists
		category, err := s.categoryRepo.FindByID(req.CategoryID)
		if err != nil {
			s.log.Error("Failed to find category", zap.Uint("category_id", req.CategoryID), zap.Error(err))
			return nil, err
		}
		if category == nil {
			s.log.Warn("Category not found", zap.Uint("category_id", req.CategoryID))
			return nil, utils.ErrCategoryNotFound
		}
		product.CategoryID = req.CategoryID
		product.Category = *category
	}

	// Check for duplicate name if name is being changed
	if req.Name != "" && req.Name != product.Name {
		existingProduct, err := s.productRepo.FindByNameAndCategory(req.Name, product.CategoryID)
		if err != nil {
			s.log.Error("Failed to check existing product name", zap.Error(err))
			return nil, err
		}
		if existingProduct != nil {
			s.log.Warn("Product name already exists in category",
				zap.String("name", req.Name),
				zap.Uint("category_id", product.CategoryID))
			return nil, utils.ErrProductExists
		}
		product.Name = req.Name
	}

	// Update other fields if provided
	if req.Description != "" {
		product.Description = req.Description
	}
	if req.Price > 0 {
		product.Price = req.Price
	}
	if req.Stock >= 0 {
		product.Stock = req.Stock
	}
	if req.MinStock > 0 {
		product.MinStock = req.MinStock
	}
	if req.ImageURL != "" {
		product.ImageURL = req.ImageURL
	}
	if req.Status != "" {
		product.Status = entity.ProductStatus(req.Status)
	}

	// Check if no changes
	if !hasProductChanges(product, req) {
		s.log.Warn("No changes provided for product update", zap.Uint("id", id))
		return nil, utils.ErrNoChangesProvided
	}

	err = s.productRepo.Update(product)
	if err != nil {
		s.log.Error("Failed to update product in repository", zap.Error(err))
		return nil, err
	}

	s.log.Info("Product updated successfully", zap.Uint("id", id))
	return response.ToProductResponse(product), nil
}

func (s *productService) DeleteProduct(id uint) error {
	s.log.Info("Deleting product", zap.Uint("id", id))

	// Check if product exists
	product, err := s.productRepo.FindByID(id)
	if err != nil {
		s.log.Error("Failed to get product for deletion", zap.Error(err))
		return err
	}

	if product == nil {
		s.log.Warn("Product not found for deletion", zap.Uint("id", id))
		return utils.ErrProductNotFound
	}

	// Check if product has order items
	hasOrderItems, err := s.productRepo.CheckHasOrderItems(id)
	if err != nil {
		s.log.Error("Failed to check product order items", zap.Error(err))
		return err
	}

	if hasOrderItems {
		s.log.Warn("Cannot delete product with associated order items", zap.Uint("id", id))
		return utils.ErrProductHasOrders
	}

	// Perform soft delete
	err = s.productRepo.SoftDelete(id)
	if err != nil {
		s.log.Error("Failed to delete product from repository", zap.Error(err))
		return err
	}

	s.log.Info("Product deleted successfully", zap.Uint("id", id))
	return nil
}

// Helper function to check if there are actual changes
func hasProductChanges(product *entity.Product, req request.UpdateProductRequest) bool {
	return (req.Name != "" && req.Name != product.Name) ||
		(req.Description != "" && req.Description != product.Description) ||
		(req.Price > 0 && req.Price != product.Price) ||
		(req.Stock >= 0 && req.Stock != product.Stock) ||
		(req.MinStock > 0 && req.MinStock != product.MinStock) ||
		(req.ImageURL != "" && req.ImageURL != product.ImageURL) ||
		(req.Status != "" && req.Status != string(product.Status)) ||
		(req.CategoryID != 0 && req.CategoryID != product.CategoryID)
}
