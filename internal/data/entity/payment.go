package entity

import (
	"time"

	"gorm.io/gorm"
)

type Payment struct {
	gorm.Model
	BookingID     uint      `gorm:"uniqueIndex;not null" json:"booking_id"`
	PaymentMethod string    `gorm:"type:varchar(50)" json:"payment_method"`
	TransactionID string    `gorm:"type:varchar(100)" json:"transaction_id"`
	PaidAt        time.Time `json:"paid_at"`
}
