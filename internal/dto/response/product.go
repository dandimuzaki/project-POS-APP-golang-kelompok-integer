package response

import (
	"project-POS-APP-golang-integer/internal/data/entity"
	"time"
)

// ProductResponse for product data
type ProductResponse struct {
	ID          uint              `json:"id"`
	Name        string            `json:"name"`
	Description string            `json:"description,omitempty"`
	Price       float64           `json:"price"`
	Stock       int               `json:"stock"`
	MinStock    int               `json:"min_stock"`
	Status      string            `json:"status"`
	ImageURL    string            `json:"image_url,omitempty"`
	CategoryID  uint              `json:"category_id"`
	Category    *CategoryResponse `json:"category,omitempty"`
	CreatedAt   time.Time         `json:"created_at"`
	UpdatedAt   time.Time         `json:"updated_at"`
}

// ProductListResponse for list of products with pagination
type ProductListResponse struct {
	Data       []ProductResponse `json:"data"`
	Pagination PaginationMeta    `json:"pagination"`
}

// Convert entity to response
func ToProductResponse(product *entity.Product) *ProductResponse {
	if product == nil {
		return nil
	}

	resp := &ProductResponse{
		ID:          product.ID,
		Name:        product.Name,
		Description: product.Description,
		Price:       product.Price,
		Stock:       product.Stock,
		MinStock:    product.MinStock,
		Status:      string(product.Status),
		ImageURL:    product.ImageURL,
		CategoryID:  product.CategoryID,
		CreatedAt:   product.CreatedAt,
		UpdatedAt:   product.UpdatedAt,
	}

	// Include category if loaded
	if product.Category.ID != 0 {
		resp.Category = &CategoryResponse{
			ID:   product.Category.ID,
			Name: product.Category.Name,
		}
	}

	return resp
}

// Convert entity slice to response slice
func ToProductListResponse(products []entity.Product) []ProductResponse {
	responses := make([]ProductResponse, len(products))
	for i, product := range products {
		responses[i] = *ToProductResponse(&product)
	}
	return responses
}
