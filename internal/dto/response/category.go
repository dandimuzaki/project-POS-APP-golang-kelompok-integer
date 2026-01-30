package response

import "time"

// CategoryResponse for list categories
type CategoryResponse struct {
	ID          uint      `json:"id"`
	Name        string    `json:"name"`
	IconURL     string    `json:"icon_url,omitempty"`
	Description string    `json:"description,omitempty"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// SuccessResponse structure
type SuccessResponse struct {
	Status  string      `json:"status"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

// ErrorResponse structure
type ErrorResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}

// PaginationResponse (Optional)
type PaginationResponse struct {
	Status     string      `json:"status"`
	Message    string      `json:"message"`
	Data       interface{} `json:"data"`
	Page       int         `json:"page"`
	Limit      int         `json:"limit"`
	TotalItems int64       `json:"total_items"`
	TotalPages int         `json:"total_pages"`
}
