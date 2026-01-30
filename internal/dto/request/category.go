package request

type GetCategoriesRequest struct {
	Name    string `form:"name" validate:"omitempty,min=1,max=100"`
	Page    int    `form:"page" validate:"min=1"`
	PerPage int    `form:"per_page" validate:"min=1,max=100"`
}

type CreateCategoryRequest struct {
	Name        string `json:"name" validate:"required,min=3,max=100"`
	IconURL     string `json:"icon_url" validate:"omitempty,url"`
	Description string `json:"description" validate:"omitempty,max=500"`
}

type UpdateCategoryRequest struct {
	Name        string `json:"name" validate:"omitempty,min=3,max=100"`
	IconURL     string `json:"icon_url" validate:"omitempty,url"`
	Description string `json:"description" validate:"omitempty,max=500"`
}
