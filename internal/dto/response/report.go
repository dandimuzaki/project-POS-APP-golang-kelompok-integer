package response

type SalesResponse struct {
	Summary SalesSummary   `json:"summary"`
	PerHour []PerHourSales `json:"per_hour"`
	Daily   []DailySales   `json:"daily"`
	Monthly []MonthlySales `json:"monthly"`
}

type SalesSummary struct {
	Total        int64        `json:"total"`
	AverageDaily float64      `json:"average_daily"`
	PeakHour     PerHourSales `json:"peak_hour"`
}

type DailySales struct {
	Date  string `json:"date"`
	Sales int    `json:"sales"`
}

type MonthlySales struct {
	Month string `json:"month"`
	Sales int    `json:"sales"`
}

type PerHourSales struct {
	Hour         string  `json:"hour"`
	AverageSales float64 `json:"average_sales"`
}

type RevenueResponse struct {
	Summary RevenueSummary   `json:"summary"`
	PerHour []PerHourRevenue `json:"per_hour"`
	Daily   []DailyRevenue   `json:"daily"`
	Monthly []MonthlyRevenue `json:"monthly"`
}

type RevenueSummary struct {
	Total        float64        `json:"total"`
	AverageDaily float64        `json:"average_daily"`
	PeakHour     PerHourRevenue `json:"peak_hour"`
}

type DailyRevenue struct {
	Date    string  `json:"date"`
	Revenue float64 `json:"revenue"`
}

type MonthlyRevenue struct {
	Month   string  `json:"month"`
	Revenue float64 `json:"revenue"`
}

type PerHourRevenue struct {
	Hour           string  `json:"hour"`
	AverageRevenue float64 `json:"average_revenue"`
}