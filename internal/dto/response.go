package dto

import (
	"project-POS-APP-golang-integer/internal/data/entity"

	"github.com/google/uuid"
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
	TotalRecords int64 `json:"total_records"`
}

type AuthResponse struct {
	Token uuid.UUID `json:"token"`
}