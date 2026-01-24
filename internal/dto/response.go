package dto

type UserResponse struct {
	Name  string `json:"name"`
	Email string `json:"email"`
	Role  string `json:"role"`
}

type Pagination struct {
	CurrentPage  int `json:"current_page"`
	Limit        int `json:"limit"`
	TotalPages   int `json:"total_pages"`
	TotalRecords int `json:"total_records"`
}