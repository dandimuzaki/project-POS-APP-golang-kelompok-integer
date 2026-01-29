package request

type UserRequest struct {
	Email string `json:"email" validate:"email"`
	Role  string `json:"role"`
}