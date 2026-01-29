package entity

import (
	"time"

	"gorm.io/gorm"
)

type UserRole string

const (
	RoleSuperAdmin UserRole = "superadmin"
	RoleAdmin      UserRole = "admin"
	RoleStaff      UserRole = "staff"
)

type User struct {
	ID           uint           `gorm:"primaryKey" json:"id"`
	Email        string         `gorm:"uniqueIndex;not null" json:"email"`
	PasswordHash string         `gorm:"not null" json:"-"`
	Role         UserRole       `gorm:"type:varchar(20);not null" json:"role"`
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
	DeletedAt    gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`

	// Relations
	Staff         *Profile         `gorm:"foreignKey:UserID" json:"profile,omitempty"`
	Sessions      []Session      `gorm:"foreignKey:UserID" json:"-"`
	Orders        []Order        `gorm:"foreignKey:CreatedBy" json:"-"`
	InventoryLogs []InventoryLog `gorm:"foreignKey:CreatedBy" json:"-"`
	Notifications []Notification `gorm:"foreignKey:UserID" json:"-"`
	OTPs          []OTP          `gorm:"foreignKey:UserID" json:"-"`
}

func (r UserRole) IsValid() bool {
	switch r {
	case RoleSuperAdmin, RoleAdmin, RoleStaff:
		return true
	}
	return false
}
