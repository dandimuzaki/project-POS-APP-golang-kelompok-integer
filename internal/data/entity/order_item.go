package entity

import (
	"time"
)

type OrderItem struct {
	ID         uint      `gorm:"primaryKey" json:"id"`
	OrderID    uint      `gorm:"index;not null" json:"order_id"`
	ProductID  uint      `gorm:"index;not null" json:"product_id"`
	Quantity   int       `gorm:"not null;default:1" json:"quantity"`
	TotalPrice float64   `gorm:"not null" json:"total_price"`
	CreatedAt  time.Time `json:"created_at"`

	// Relations
	Order   Order   `gorm:"foreignKey:OrderID" json:"order,omitempty"`
	Product Product `gorm:"foreignKey:ProductID" json:"product"`
}
