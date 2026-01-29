package response

type UserResponse struct {
	Email    string `json:"email"`
	Role     string `json:"role"`
	Password string `json:"password"`
}