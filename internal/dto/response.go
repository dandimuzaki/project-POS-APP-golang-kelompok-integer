package dto

type ResponseUser struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}

type Pagination struct {
	CurrentPage  int   `json:"current_page"`
	Limit        int   `json:"limit"`
	TotalPages   int   `json:"total_pages"`
	TotalRecords int64 `json:"total_records"`
}

type TourResponse struct {
	ID            uint    `json:"id"`
	Name          string  `json:"name"`
	Location      string  `json:"location"`
	ImageURL      string  `json:"image_url"`
	Price         float64 `json:"price"`
	Duration      int     `json:"duration_day"`
	Rating        float64 `json:"rating"`
	ReviewCount   int     `json:"review_count"`
	RemainingSeat int     `json:"remaining_seat"`
}

type Tour struct {
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Rating      float64 `json:"rating"`
	ReviewCount int     `json:"review_count"`
}

type Price struct {
	BasePrice     float64 `json:"base_price"`
	PriceOverride float64 `json:"price_override"`
}

type ScheduleResponse struct {
	ID          uint     `json:"id"`
	Tour        Tour     `json:"tour"`
	Price       Price    `json:"price"`
	Destination string   `json:"destination"`
	StartDate   string   `json:"start_date"`
	EndDate     string   `json:"end_date"`
	Images      []string `json:"images"`
}