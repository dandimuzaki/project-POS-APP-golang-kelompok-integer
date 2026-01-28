package response

type PaginatedResponse[T any] struct {
	Data       []T            `json:"data"`
	Pagination PaginationMeta `json:"pagination"`
}

type PaginationMeta struct {
	Page       int   `json:"page"`
	PerPage    int   `json:"per_page"`
	Total      int64 `json:"total"`
	TotalPages int   `json:"total_pages"`
}

func NewPaginatedResponse[T any](
	data []T,
	page int,
	perPage int,
	total int64,
) *PaginatedResponse[T] {

	totalPages := 0
	if perPage > 0 && total > 0 {
		totalPages = int((total + int64(perPage) - 1) / int64(perPage))
	}

	return &PaginatedResponse[T]{
		Data: data,
		Pagination: PaginationMeta{
			Page:       page,
			PerPage:    perPage,
			Total:      total,
			TotalPages: totalPages,
		},
	}
}
