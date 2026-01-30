package request

type PaginationRequest struct {
	Page    int `json:"page" form:"page" query:"page"`
	PerPage int `json:"per_page" form:"per_page" query:"per_page"`
}

// Getter dengan default values
func (p PaginationRequest) GetPage() int {
	if p.Page < 1 {
		return 1
	}
	return p.Page
}

func (p PaginationRequest) GetPerPage() int {
	if p.PerPage < 1 {
		return 10 // DEFAULT 10
	}
	if p.PerPage > 100 {
		return 100 // MAX 100
	}
	return p.PerPage
}

func (p PaginationRequest) GetOffset() int {
	return (p.GetPage() - 1) * p.GetPerPage()
}