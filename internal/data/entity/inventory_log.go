package entity

import (
	"time"
)

// InventoryLogType enum
type InventoryLogType string

const (
	InventoryLogTypeIn         InventoryLogType = "in"
	InventoryLogTypeOut        InventoryLogType = "out"
	InventoryLogTypeAdjustment InventoryLogType = "adjustment"
	InventoryLogTypeInitial    InventoryLogType = "initial"
)

type InventoryLog struct {
	ID                uint             `gorm:"primaryKey" json:"id"`
	ProductID         uint             `gorm:"index;not null" json:"product_id"`
	Type              InventoryLogType `gorm:"type:varchar(20);not null" json:"type"`
	QuantityChange    int              `gorm:"not null" json:"quantity_change"`
	CurrentStockAfter int              `gorm:"not null" json:"current_stock_after"`
	ReferenceID       *uint            `gorm:"index" json:"reference_id,omitempty"`
	ReferenceType     string           `gorm:"type:varchar(50)" json:"reference_type,omitempty"`
	Notes             string           `json:"notes,omitempty"`
	CreatedBy         uint             `gorm:"index;not null" json:"created_by"`
	CreatedAt         time.Time        `json:"created_at"`

	// Relations
	Product Product `gorm:"foreignKey:ProductID" json:"product"`
	Creator User    `gorm:"foreignKey:CreatedBy" json:"creator"`
}
