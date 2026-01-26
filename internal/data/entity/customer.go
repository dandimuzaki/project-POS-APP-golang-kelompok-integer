package entity

import (
	"time"

	"gorm.io/gorm"
)

// CustomerTitle enum
type CustomerTitle string

const (
	CustomerTitleMr   CustomerTitle = "Mr"
	CustomerTitleMrs  CustomerTitle = "Mrs"
	CustomerTitleMs   CustomerTitle = "Ms"
	CustomerTitleDr   CustomerTitle = "Dr"
	CustomerTitleProf CustomerTitle = "Prof"
)

type Customer struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	Title     CustomerTitle  `gorm:"type:varchar(10)" json:"title,omitempty"`
	FirstName string         `gorm:"not null" json:"first_name"`
	LastName  string         `json:"last_name"`
	Phone     string         `json:"phone"`
	Email     string         `json:"email"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`

	// Relations
	Orders       []Order       `gorm:"foreignKey:CustomerID" json:"-"`
	Reservations []Reservation `gorm:"foreignKey:CustomerID" json:"-"`
}
