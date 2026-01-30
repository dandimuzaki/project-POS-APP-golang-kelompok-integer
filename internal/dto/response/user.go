package response

import (
	"project-POS-APP-golang-integer/internal/data/entity"
)

type UserResponse struct {
	Email string          `json:"email"`
	Role  entity.UserRole `json:"role"`
}

type CreateUserResponse struct {
	Email string          `json:"email"`
	Role  entity.UserRole `json:"role"`
	Password string `json:"password"`
}