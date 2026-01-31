package request

import "github.com/go-playground/validator/v10"

var validate = validator.New()

// CreateProductRequest for creating new product
type CreateProductRequest struct {
	Name        string  `json:"name" validate:"required,min=3,max=100"`
	Description string  `json:"description" validate:"omitempty,max=500"`
	Price       float64 `json:"price" validate:"required,min=0"`
	Stock       int     `json:"stock" validate:"min=0"`
	MinStock    int     `json:"min_stock" validate:"min=0"`
	ImageURL    string  `json:"image_url" validate:"omitempty,url"`
	CategoryID  uint    `json:"category_id" validate:"required,min=1"`
	Status      string  `json:"status" validate:"omitempty,oneof=active inactive"`
}

// UpdateProductRequest for updating existing product
type UpdateProductRequest struct {
	Name        string  `json:"name" validate:"omitempty,min=3,max=100"`
	Description string  `json:"description" validate:"omitempty,max=500"`
	Price       float64 `json:"price" validate:"omitempty,min=0"`
	Stock       int     `json:"stock" validate:"omitempty,min=0"`
	MinStock    int     `json:"min_stock" validate:"omitempty,min=0"`
	ImageURL    string  `json:"image_url" validate:"omitempty,url"`
	CategoryID  uint    `json:"category_id" validate:"omitempty,min=1"`
	Status      string  `json:"status" validate:"omitempty,oneof=active inactive"`
}

// GetProductsRequest for filtering and sorting products
type GetProductsRequest struct {
	Page         int     `form:"page" validate:"min=1"`
	Limit        int     `form:"limit" validate:"min=1,max=100"`
	CategoryID   uint    `form:"category_id" validate:"omitempty,min=1"`
	CategoryName string  `form:"category_name" validate:"omitempty"`
	Status       string  `form:"status" validate:"omitempty,oneof=active inactive"`
	MinPrice     float64 `form:"min_price" validate:"omitempty,min=0"`
	MaxPrice     float64 `form:"max_price" validate:"omitempty,min=0"`
	Search       string  `form:"search" validate:"omitempty"`

	// Sorting
	SortByStock     string `form:"sort_by_stock" validate:"omitempty,oneof=asc desc"`
	SortByPrice     string `form:"sort_by_price" validate:"omitempty,oneof=asc desc"`
	SortByCreatedAt string `form:"sort_by_created_at" validate:"omitempty,oneof=asc desc"`
	SortBySold      string `form:"sort_by_sold" validate:"omitempty,oneof=asc desc"`
}

// UpdateStockRequest for updating product stock
type UpdateStockRequest struct {
	Quantity int    `json:"quantity" validate:"required"`
	Note     string `json:"note" validate:"omitempty,max=500"`
}

// Validate methods
func (req *CreateProductRequest) Validate() error {
	return validate.Struct(req)
}

func (req *UpdateProductRequest) Validate() error {
	return validate.Struct(req)
}

func (req *GetProductsRequest) Validate() error {
	return validate.Struct(req)
}

func (req *UpdateStockRequest) Validate() error {
	return validate.Struct(req)
}
