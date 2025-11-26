package models

import (
	"time"

	"gorm.io/gorm"
)

type NotificationLog struct {
	ID             uint      `gorm:"primaryKey" json:"id"`
	SubscriptionID uint      `gorm:"not null;index:idx_subscription_sent" json:"subscription_id"`
	SentAt         time.Time `gorm:"not null;index:idx_subscription_sent" json:"sent_at"`
	Status         string    `gorm:"not null" json:"status"` // "success", "failed"
	ErrorMessage   string    `json:"error_message,omitempty"`

	// Relationships
	Subscription Subscription `gorm:"foreignKey:SubscriptionID" json:"subscription,omitempty"`
}

// BeforeCreate is a GORM hook that runs before creating a notification log
func (n *NotificationLog) BeforeCreate(tx *gorm.DB) error {
	n.SentAt = time.Now()
	return nil
}
