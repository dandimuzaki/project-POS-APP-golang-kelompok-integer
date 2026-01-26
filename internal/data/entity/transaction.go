package entity

import (
	"time"
)

// TransactionType enum
type TransactionType string

const (
	TransactionTypePayment    TransactionType = "payment"
	TransactionTypeRefund     TransactionType = "refund"
	TransactionTypeAdjustment TransactionType = "adjustment"
)

// TransactionStatus enum
type TransactionStatus string

const (
	TransactionStatusPending   TransactionStatus = "pending"
	TransactionStatusCompleted TransactionStatus = "completed"
	TransactionStatusFailed    TransactionStatus = "failed"
	TransactionStatusCancelled TransactionStatus = "cancelled"
)

type Transaction struct {
	ID                uint              `gorm:"primaryKey" json:"id"`
	TransactionNumber string            `gorm:"uniqueIndex;not null" json:"transaction_number"`
	OrderID           uint              `gorm:"index;not null" json:"order_id"`
	TransactionType   TransactionType   `gorm:"type:varchar(20);not null" json:"transaction_type"`
	PaymentMethodID   uint              `gorm:"index;not null" json:"payment_method_id"`
	Amount            float64           `gorm:"not null" json:"amount"`
	Status            TransactionStatus `gorm:"type:varchar(20);default:'pending'" json:"status"`
	Notes             string            `json:"notes,omitempty"`
	CreatedBy         uint              `gorm:"index;not null" json:"created_by"`
	CreatedAt         time.Time         `json:"created_at"`
	UpdatedAt         time.Time         `json:"updated_at"`

	// Relations
	Order         Order         `gorm:"foreignKey:OrderID" json:"order"`
	PaymentMethod PaymentMethod `gorm:"foreignKey:PaymentMethodID" json:"payment_method"`
	Creator       User          `gorm:"foreignKey:CreatedBy" json:"creator"`
}
