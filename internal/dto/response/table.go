package response

import (
	"project-POS-APP-golang-integer/internal/data/entity"
	"time"
)

type TableResponse struct {
	ID          uint               `json:"id"`
	TableNumber string             `json:"table_number"`
	Capacity    int                `json:"capacity"`
	Status      entity.TableStatus `json:"status"`
	CreatedAt   time.Time          `json:"created_at"`
	UpdatedAt   time.Time          `json:"updated_at"`
}

// Converters
func TableToResponse(table *entity.Table) TableResponse {
	return TableResponse{
		ID:          table.ID,
		TableNumber: table.TableNumber,
		Capacity:    table.Capacity,
		Status:      table.Status,
		CreatedAt:   table.CreatedAt,
		UpdatedAt:   table.UpdatedAt,
	}
}
