package request

type UserFilterRequest struct {
	PaginationRequest
	Role  string `form:"role"`
	Name  string `form:"name"`
	Email string `form:"email"`
}

type CreateUserRequest struct {
	Email             string  `json:"email" validate:"email"`
	Role              string  `json:"role"`
	FullName          string  `json:"full_name"`
	Phone             string  `json:"phone"`
	DateOfBirth       string  `json:"date_of_birth"`
	Salary            float64 `json:"salary"`
	ProfileImageURL   string  `json:"profile_image_url"`
	Address           string  `json:"address"`
	AdditionalDetails string  `json:"additional_details,omitempty"`
}

type UpdateUserRequest struct {
	Email             string  `json:"email" validate:"email"`
	Role              string  `json:"role"`
	FullName          string  `json:"full_name"`
	Phone             string  `json:"phone"`
	DateOfBirth       string  `json:"date_of_birth"`
	Salary            float64 `json:"salary"`
	ProfileImageURL   string  `json:"profile_image_url"`
	Address           string  `json:"address"`
	AdditionalDetails string  `json:"additional_details,omitempty"`
}