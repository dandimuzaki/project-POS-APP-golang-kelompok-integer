package entity

import "gorm.io/gorm"

type Review struct {
	gorm.Model
	TourID uint `gorm:"not null" json:"tour_id"`
	UserID uint `gorm:"not null" json:"user_id"`
	Rating int `gorm:"check:rating >=1 AND rating <= 5" json:"rating"`
	Comment string `gorm:"type:text" json:"comment"`
	User User `gorm:"foreignKey:UserID" json:"user"`
}