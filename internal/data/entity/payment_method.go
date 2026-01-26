package entity

import (
	"time"
)

type PaymentMethod struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Name      string    `gorm:"uniqueIndex;not null" json:"name"`
	IsActive  bool      `gorm:"default:true" json:"is_active"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	// Relations
	Orders       []Order       `gorm:"foreignKey:PaymentMethodID" json:"-"`
	Transactions []Transaction `gorm:"foreignKey:PaymentMethodID" json:"-"`
}
