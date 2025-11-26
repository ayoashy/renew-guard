package repositories

import (
	"renew-guard/internal/models"

	"gorm.io/gorm"
)

type NotificationLogRepository interface {
	Create(log *models.NotificationLog) error
	FindBySubscriptionID(subscriptionID uint) ([]models.NotificationLog, error)
}

type notificationLogRepository struct {
	db *gorm.DB
}

func NewNotificationLogRepository(db *gorm.DB) NotificationLogRepository {
	return &notificationLogRepository{db: db}
}

func (r *notificationLogRepository) Create(log *models.NotificationLog) error {
	return r.db.Create(log).Error
}

func (r *notificationLogRepository) FindBySubscriptionID(subscriptionID uint) ([]models.NotificationLog, error) {
	var logs []models.NotificationLog
	err := r.db.Where("subscription_id = ?", subscriptionID).
		Order("sent_at DESC").
		Find(&logs).Error
	return logs, err
}
