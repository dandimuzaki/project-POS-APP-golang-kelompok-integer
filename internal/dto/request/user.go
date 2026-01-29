package request

type UserFilterRequest struct {
	PaginationRequest
	Role  string `form:"role"`
	Name  string `form:"name"`
	Email string `form:"email"`
}

type UserRequest struct {
	Email string `json:"email" validate:"email"`
	Role  string `json:"role"`
}

type UpdateUserRequest struct {
	ID   uint   `json:"id"`
	Role string `json:"role"`
}