package response

import (
	"project-POS-APP-golang-integer/internal/data/entity"
	"time"
)

type InventoryLogResponse struct {
	ID                uint                    `json:"id"`
	ProductID         uint                    `json:"product_id"`
	Type              entity.InventoryLogType `json:"type"`
	QuantityChange    int                     `json:"quantity_change"`
	CurrentStockAfter int                     `json:"current_stock_after"`
	ReferenceID       *uint                   `json:"reference_id,omitempty"`
	ReferenceType     string                  `json:"reference_type,omitempty"`
	Notes             string                  `json:"notes,omitempty"`
	CreatedBy         uint                    `json:"created_by"`
	CreatedAt         time.Time               `json:"created_at"`
}