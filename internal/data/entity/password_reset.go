package entity

import (
	"time"
)

type PasswordReset struct {
	ID        uint `gorm:"primaryKey" json:"id"`
	UserID uint `gorm:"not null;index" json:"user_id"`
	ResetTokenHash string `gorm:"not null;uniqueIndex" json:"reset_token_hash"`
	ExpiredAt time.Time `gorm:"not null" json:"expired_at"`
	UsedAt    *time.Time `json:"used_at"`
	CreatedAt time.Time `json:"created_at"`

	// Relations
	User User `gorm:"foreignKey:UserID" json:"user"`
}
