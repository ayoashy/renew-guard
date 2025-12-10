package models

import (
	"time"

	"gorm.io/gorm"
)

type Subscription struct {
	ID                   uint       `gorm:"primaryKey" json:"id"`
	UserID               uint       `gorm:"not null;index" json:"user_id"`
	Email                string     `gorm:"default:''" json:"email"` // User's email at subscription creation
	Name                 string     `gorm:"not null" json:"name"`
	StartDate            time.Time  `gorm:"not null" json:"start_date"`
	DurationDays         int        `gorm:"not null" json:"duration_days"`
	EndDate              time.Time  `gorm:"not null;index" json:"end_date"`
	NotificationEnabled  bool       `gorm:"default:true" json:"notification_enabled"`
	LastNotificationSent *time.Time `json:"last_notification_sent,omitempty"`
	CreatedAt            time.Time  `json:"created_at"`
	UpdatedAt            time.Time  `json:"updated_at"`

	// Relationships
	User             User              `gorm:"foreignKey:UserID" json:"user,omitempty"`
	NotificationLogs []NotificationLog `gorm:"foreignKey:SubscriptionID" json:"notification_logs,omitempty"`
}

// BeforeCreate is a GORM hook that runs before creating a subscription
func (s *Subscription) BeforeCreate(tx *gorm.DB) error {
	s.CreatedAt = time.Now()
	s.UpdatedAt = time.Now()
	s.ComputeEndDate()
	return nil
}

// BeforeUpdate is a GORM hook that runs before updating a subscription
func (s *Subscription) BeforeUpdate(tx *gorm.DB) error {
	s.UpdatedAt = time.Now()
	s.ComputeEndDate()
	return nil
}

// ComputeEndDate calculates the end date based on start date and duration
func (s *Subscription) ComputeEndDate() {
	s.EndDate = s.StartDate.AddDate(0, 0, s.DurationDays)
}

// DaysUntilExpiration returns the number of days until the subscription expires
func (s *Subscription) DaysUntilExpiration() int {
	now := time.Now()
	duration := s.EndDate.Sub(now)
	days := int(duration.Hours() / 24)
	return days
}

// IsExpired checks if the subscription has already expired
func (s *Subscription) IsExpired() bool {
	return time.Now().After(s.EndDate)
}

// ShouldNotify determines if a notification should be sent
func (s *Subscription) ShouldNotify(daysBefore int) bool {
	if !s.NotificationEnabled {
		return false
	}

	if s.IsExpired() {
		return false
	}

	daysLeft := s.DaysUntilExpiration()
	if daysLeft > daysBefore || daysLeft < 0 {
		return false
	}

	// Check if notification was already sent today
	if s.LastNotificationSent != nil {
		today := time.Now().Truncate(24 * time.Hour)
		lastSent := s.LastNotificationSent.Truncate(24 * time.Hour)
		if today.Equal(lastSent) {
			return false
		}
	}

	return true
}
