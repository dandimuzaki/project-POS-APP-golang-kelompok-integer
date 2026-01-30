package request

type CreateInventoryLogRequest struct {
	ProductID      uint   `json:"product_id"`
	Action         string `json:"action"`
	QuantityChange int    `json:"quantity_change"`
	Notes          string `json:"notes"`
}

type InventoryLogsFilter struct {
	PaginationRequest
}