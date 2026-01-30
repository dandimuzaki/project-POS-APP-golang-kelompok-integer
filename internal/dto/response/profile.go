package response

import "project-POS-APP-golang-integer/internal/data/entity"

type ProfileResponse struct {
	Email             string          `json:"email" validate:"email"`
	Role              entity.UserRole `json:"role"`
	FullName          string          `gorm:"not null" json:"full_name"`
	Phone             string          `json:"phone"`
	DateOfBirth       string          `json:"date_of_birth"`
	Salary            float64         `json:"salary"`
	ProfileImageURL   string          `json:"profile_image_url"`
	Address           string          `json:"address"`
	AdditionalDetails string          `gorm:"type:text" json:"additional_details,omitempty"`
}