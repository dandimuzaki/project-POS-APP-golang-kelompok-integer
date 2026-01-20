package entity

import (
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Name     string    `gorm:"type:varchar(100);not null" json:"name"`
	Email    string    `gorm:"type:varchar(100);uniqueIndex;not null" json:"email"`
	Password string    `gorm:"type:varchar(255);not null" json:"-"`
	Role     string    `gorm:"type:varchar(20);default:'customer'" json:"role"`
	Bookings []Booking `gorm:"foreignKey:UserID" json:"bookings,omitempty"`
	Reviews  []Review  `gorm:"foreignKey:UserID" json:"reviews,omitempty"`
}
