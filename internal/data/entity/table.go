package entity

import (
	"time"
)

// TableStatus enum
type TableStatus string

const (
	TableStatusAvailable TableStatus = "available"
	TableStatusOccupied  TableStatus = "occupied"
	TableStatusReserved  TableStatus = "reserved"
)

type Table struct {
	ID          uint        `gorm:"primaryKey" json:"id"`
	TableNumber string      `gorm:"uniqueIndex;not null" json:"table_number"`
	Capacity    int         `gorm:"not null" json:"capacity"`
	Status      TableStatus `gorm:"type:varchar(20);default:'available'" json:"status"`
	CreatedAt   time.Time   `json:"created_at"`
	UpdatedAt   time.Time   `json:"updated_at"`

	// Relations
	Orders       []Order       `gorm:"foreignKey:TableID" json:"-"`
	Reservations []Reservation `gorm:"foreignKey:TableID" json:"-"`
}
