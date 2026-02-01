package request

import "time"

type PeriodRequest struct {
	From string `json:"from"`
	To   string `json:"to"`
}

type PeriodQuery struct {
	From time.Time `json:"from"`
	To   time.Time `json:"to"`
}