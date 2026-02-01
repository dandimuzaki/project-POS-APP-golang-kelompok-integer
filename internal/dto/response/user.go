package response

import (
	"project-POS-APP-golang-integer/internal/data/entity"
)

type UserResponse struct {
	Email string          `json:"email"`
	Role  entity.UserRole `json:"role"`
	FullName          string  `json:"full_name"`
	Phone             string  `json:"phone"`
	DateOfBirth       string  `json:"date_of_birth"`
	Salary            float64 `json:"salary"`
	ProfileImageURL   string  `json:"profile_image_url"`
	Address           string  `json:"address"`
	AdditionalDetails string  `json:"additional_details,omitempty"`
}

type CreateUserResponse struct {
	Email string          `json:"email"`
	Role  entity.UserRole `json:"role"`
	Password string `json:"password"`
	FullName          string  `json:"full_name"`
	Phone             string  `json:"phone"`
	DateOfBirth       string  `json:"date_of_birth"`
	Salary            float64 `json:"salary"`
	ProfileImageURL   string  `json:"profile_image_url"`
	Address           string  `json:"address"`
	AdditionalDetails string  `json:"additional_details,omitempty"`
}