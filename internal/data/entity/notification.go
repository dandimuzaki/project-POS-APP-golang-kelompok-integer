package entity

import (
	"time"

	"gorm.io/datatypes"
)

// NotificationType enum
type NotificationType string

const (
	NotificationTypeStockAlert     NotificationType = "stock_alert"
	NotificationTypeNewOrder       NotificationType = "new_order"
	NotificationTypeNewReservation NotificationType = "new_reservation"
	NotificationTypeSystem         NotificationType = "system"
)

// NotificationStatus enum
type NotificationStatus string

const (
	NotificationStatusNew  NotificationStatus = "new"
	NotificationStatusRead NotificationStatus = "read"
)

type Notification struct {
	ID        uint               `gorm:"primaryKey" json:"id"`
	UserID    uint               `gorm:"index;not null" json:"user_id"`
	Title     string             `gorm:"not null" json:"title"`
	Message   string             `gorm:"type:text;not null" json:"message"`
	Type      NotificationType   `gorm:"type:varchar(30);not null" json:"type"`
	Status    NotificationStatus `gorm:"type:varchar(20);default:'new'" json:"status"`
	Metadata  datatypes.JSON     `json:"metadata,omitempty"`
	ReadAt    *time.Time         `json:"read_at,omitempty"`
	CreatedAt time.Time          `json:"created_at"`
	UpdatedAt time.Time          `json:"updated_at"`

	// Relations
	User User `gorm:"foreignKey:UserID" json:"user"`
}
