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

// Validate validates the request struct
func (req *CreateProductRequest) Validate() error {
	return validate.Struct(req)
}

func (req *UpdateProductRequest) Validate() error {
	return validate.Struct(req)
}
