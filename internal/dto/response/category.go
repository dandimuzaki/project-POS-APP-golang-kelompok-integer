package response

import "time"

type CategoryResponse struct {
	ID          uint      `json:"id"`
	Name        string    `json:"name"`
	IconURL     string    `json:"icon_url,omitempty"`
	Description string    `json:"description,omitempty"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type CategoryListResponse struct {
	Data       []CategoryResponse `json:"data"`
	Pagination PaginationMeta     `json:"pagination"`
}
