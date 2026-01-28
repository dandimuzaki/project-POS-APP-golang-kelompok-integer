package request

type CreateTableRequest struct {
	TableNumber string `json:"table_number" form:"table_number" validate:"required"`
	Capacity    int    `json:"capacity" form:"capacity" validate:"required,min=1,max=20"`
}

type UpdateTableRequest struct {
	TableNumber string `json:"table_number" form:"table_number"`
	Capacity    int    `json:"capacity" form:"capacity" validate:"omitempty,min=1,max=20"`
	Status      string `json:"status" form:"status" validate:"omitempty,oneof=available occupied reserved"`
}

type GetTablesRequest struct {
	PaginationRequest
	Status      string `json:"status" form:"status"`
	MinCapacity int    `json:"min_capacity" form:"min_capacity" validate:"omitempty,min=1"`
}
