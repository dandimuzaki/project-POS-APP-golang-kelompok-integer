package dto

import (
	"project-POS-APP-golang-integer/internal/data/entity"
	"time"
)

type UserResponse struct {
	Name  string          `json:"name"`
	Email string          `json:"email"`
	Role  entity.UserRole `json:"role"`
}

type Pagination struct {
	CurrentPage  int `json:"current_page"`
	Limit        int `json:"limit"`
	TotalPages   int `json:"total_pages"`
	TotalRecords int `json:"total_records"`
}

type CategoryResponse struct {
	ID          uint      `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}
