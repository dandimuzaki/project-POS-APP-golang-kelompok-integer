package dto

type TourFilterRequest struct {
	Page   int    `form:"page" binding:"min=1"`
	Limit  int    `form:"limit" binding:"min=1,max=100"`
	Search string `form:"search"`
	Date   string `form:"date"`
	SortBy string `form:"sort_by"`
}