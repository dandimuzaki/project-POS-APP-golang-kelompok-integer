package entity

import (
	"time"

	"gorm.io/gorm"
)

// ReservationStatus enum
type ReservationStatus string

const (
	ReservationStatusAwaiting  ReservationStatus = "awaiting"
	ReservationStatusConfirmed ReservationStatus = "confirmed"
	ReservationStatusCancelled ReservationStatus = "cancelled"
	ReservationStatusCompleted ReservationStatus = "completed"
)

type Reservation struct {
	ID              uint              `gorm:"primaryKey" json:"id"`
	CustomerID      uint              `gorm:"index;not null" json:"customer_id"`
	TableID         uint              `gorm:"index;not null" json:"table_id"`
	PaxNumber       int               `gorm:"not null" json:"pax_number"`
	ReservationDate time.Time         `gorm:"not null" json:"reservation_date"`
	ReservationTime time.Time         `gorm:"not null" json:"reservation_time"`
	DepositFee      float64           `gorm:"default:0" json:"deposit_fee"`
	Status          ReservationStatus `gorm:"type:varchar(20);default:'awaiting'" json:"status"`
	Notes           string            `json:"notes,omitempty"`
	CheckOutAt      *time.Time        `json:"check_out_at,omitempty"`
	CreatedAt       time.Time         `json:"created_at"`
	UpdatedAt       time.Time         `json:"updated_at"`
	DeletedAt       gorm.DeletedAt    `gorm:"index" json:"deleted_at,omitempty"`

	// Relations
	Customer Customer `gorm:"foreignKey:CustomerID" json:"customer"`
	Table    Table    `gorm:"foreignKey:TableID" json:"table"`
}
