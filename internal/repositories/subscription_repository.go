package repositories

import (
	"renew-guard/internal/models"
	"time"

	"gorm.io/gorm"
)

type SubscriptionRepository interface {
	Create(subscription *models.Subscription) error
	FindByID(id uint) (*models.Subscription, error)
	FindByUserID(userID uint) ([]models.Subscription, error)
	Update(subscription *models.Subscription) error
	Delete(id uint) error
	FindExpiringSubscriptions(daysBefore int) ([]models.Subscription, error)
	UpdateLastNotificationSent(id uint, sentAt time.Time) error
}

type subscriptionRepository struct {
	db *gorm.DB
}

func NewSubscriptionRepository(db *gorm.DB) SubscriptionRepository {
	return &subscriptionRepository{db: db}
}

func (r *subscriptionRepository) Create(subscription *models.Subscription) error {
	return r.db.Create(subscription).Error
}

func (r *subscriptionRepository) FindByID(id uint) (*models.Subscription, error) {
	var subscription models.Subscription
	err := r.db.Preload("User").First(&subscription, id).Error
	if err != nil {
		return nil, err
	}
	return &subscription, nil
}

func (r *subscriptionRepository) FindByUserID(userID uint) ([]models.Subscription, error) {
	var subscriptions []models.Subscription
	err := r.db.Where("user_id = ?", userID).Order("end_date ASC").Find(&subscriptions).Error
	return subscriptions, err
}

func (r *subscriptionRepository) Update(subscription *models.Subscription) error {
	return r.db.Save(subscription).Error
}

func (r *subscriptionRepository) Delete(id uint) error {
	return r.db.Delete(&models.Subscription{}, id).Error
}

func (r *subscriptionRepository) FindExpiringSubscriptions(daysBefore int) ([]models.Subscription, error) {
	var subscriptions []models.Subscription
	
	now := time.Now()
	targetDate := now.AddDate(0, 0, daysBefore)

	err := r.db.Preload("User").
		Where("notification_enabled = ?", true).
		Where("end_date >= ?", now).
		Where("end_date <= ?", targetDate).
		Find(&subscriptions).Error

	return subscriptions, err
}

func (r *subscriptionRepository) UpdateLastNotificationSent(id uint, sentAt time.Time) error {
	return r.db.Model(&models.Subscription{}).
		Where("id = ?", id).
		Update("last_notification_sent", sentAt).Error
}
