package dto

type UserFilterRequest struct {
	Page  int    `form:"page" binding:"min=1"`
	Limit int    `form:"limit" binding:"min=1,max=100"`
	Role  string `form:"role"`
}

type CategoryRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}
