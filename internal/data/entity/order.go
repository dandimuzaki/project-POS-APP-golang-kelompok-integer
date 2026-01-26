package entity

import (
	"gorm.io/gorm"
)

// OrderStatus enum
type OrderStatus string

const (
	OrderStatusPending   OrderStatus = "pending"
	OrderStatusInProcess OrderStatus = "in_process"
	OrderStatusCooking   OrderStatus = "cooking"
	OrderStatusCompleted OrderStatus = "completed"
	OrderStatusCancelled OrderStatus = "cancelled"
)

type Order struct {
	gorm.Model
	OrderNumber     string      `gorm:"uniqueIndex;not null" json:"order_number"`
	CustomerID      *uint       `gorm:"index" json:"customer_id,omitempty"`
	TableID         uint        `gorm:"index;not null" json:"table_id"`
	Status          OrderStatus `gorm:"type:varchar(20);default:'pending'" json:"status"`
	StatusDesc      string      `gorm:"type:varchar(100)" json:"status_desc,omitempty"`
	Subtotal        float64     `gorm:"not null;default:0" json:"subtotal"`
	TaxPercentage   float64     `gorm:"not null;default:10" json:"tax_percentage"`
	TaxAmount       float64     `gorm:"not null;default:0" json:"tax_amount"`
	Total           float64     `gorm:"not null;default:0" json:"total"`
	PaymentMethodID uint        `gorm:"index;not null" json:"payment_method_id"`
	CreatedBy       uint        `gorm:"index;not null" json:"created_by"`
	Notes           string      `json:"notes,omitempty"`

	// Relations
	Customer      Customer      `gorm:"foreignKey:CustomerID" json:"customer,omitempty"`
	Table         Table         `gorm:"foreignKey:TableID" json:"table"`
	PaymentMethod PaymentMethod `gorm:"foreignKey:PaymentMethodID" json:"payment_method"`
	Creator       User          `gorm:"foreignKey:CreatedBy" json:"creator"`
	OrderItems    []OrderItem   `gorm:"foreignKey:OrderID" json:"items"`
	Transactions  []Transaction `gorm:"foreignKey:OrderID" json:"transactions,omitempty"`
}
