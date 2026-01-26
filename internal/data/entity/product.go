package entity

import (
	"gorm.io/gorm"
)

type ProductStatus string

const (
	ProductStatusActive   ProductStatus = "active"
	ProductStatusInactive ProductStatus = "inactive"
)

type Product struct {
	gorm.Model
	Name        string        `gorm:"not null" json:"name"`
	Description string        `json:"description,omitempty"`
	Price       float64       `gorm:"not null;default:0" json:"price"`
	Stock       int           `gorm:"default:0" json:"stock"`
	MinStock    int           `gorm:"default:5" json:"min_stock"`
	Status      ProductStatus `gorm:"type:varchar(20);default:'active'" json:"status"`
	ImageURL    string        `json:"image_url,omitempty"`
	CategoryID  uint          `gorm:"index;not null" json:"category_id"`

	// Relations
	Category      Category       `gorm:"foreignKey:CategoryID" json:"category,omitempty"`
	OrderItems    []OrderItem    `gorm:"foreignKey:ProductID" json:"-"`
	InventoryLogs []InventoryLog `gorm:"foreignKey:ProductID" json:"-"`
}
