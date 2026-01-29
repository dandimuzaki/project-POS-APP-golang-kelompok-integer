package dto

type UserFilterRequest struct {
	Page  int    `form:"page" binding:"min=1"`
	Limit int    `form:"limit" binding:"min=1,max=100"`
	Role  string `form:"role"`
}

type CreateUserRequest struct {
	Email string `json:"email" validate:"email"`
	Role  string `json:"role"`
}

type LoginRequest struct {
	Email    string `json:"email" validate:"email"`
	Password string `json:"password"`
}