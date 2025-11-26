package services

import (
	"errors"
	"renew-guard/internal/models"
	"renew-guard/internal/repositories"
	"time"

	"gorm.io/gorm"
)

var (
	ErrSubscriptionNotFound   = errors.New("subscription not found")
	ErrUnauthorizedAccess     = errors.New("unauthorized access to subscription")
	ErrInvalidSubscriptionData = errors.New("invalid subscription data")
)

type SubscriptionService interface {
	Create(userID uint, name string, startDate time.Time, durationDays int) (*models.Subscription, error)
	GetByID(id, userID uint) (*models.Subscription, error)
	GetAllByUserID(userID uint) ([]models.Subscription, error)
	Update(id, userID uint, name string, startDate time.Time, durationDays int, notificationEnabled bool) (*models.Subscription, error)
	Delete(id, userID uint) error
	ToggleNotification(id, userID uint, enabled bool) (*models.Subscription, error)
}

type subscriptionService struct {
	subscriptionRepo repositories.SubscriptionRepository
}

func NewSubscriptionService(subscriptionRepo repositories.SubscriptionRepository) SubscriptionService {
	return &subscriptionService{
		subscriptionRepo: subscriptionRepo,
	}
}

func (s *subscriptionService) Create(userID uint, name string, startDate time.Time, durationDays int) (*models.Subscription, error) {
	// Validate input
	if name == "" || durationDays <= 0 {
		return nil, ErrInvalidSubscriptionData
	}

	subscription := &models.Subscription{
		UserID:              userID,
		Name:                name,
		StartDate:           startDate,
		DurationDays:        durationDays,
		NotificationEnabled: true,
	}

	if err := s.subscriptionRepo.Create(subscription); err != nil {
		return nil, err
	}

	return subscription, nil
}

func (s *subscriptionService) GetByID(id, userID uint) (*models.Subscription, error) {
	subscription, err := s.subscriptionRepo.FindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrSubscriptionNotFound
		}
		return nil, err
	}

	// Verify ownership
	if subscription.UserID != userID {
		return nil, ErrUnauthorizedAccess
	}

	return subscription, nil
}

func (s *subscriptionService) GetAllByUserID(userID uint) ([]models.Subscription, error) {
	return s.subscriptionRepo.FindByUserID(userID)
}

func (s *subscriptionService) Update(id, userID uint, name string, startDate time.Time, durationDays int, notificationEnabled bool) (*models.Subscription, error) {
	// Get existing subscription and verify ownership
	subscription, err := s.GetByID(id, userID)
	if err != nil {
		return nil, err
	}

	// Validate input
	if name == "" || durationDays <= 0 {
		return nil, ErrInvalidSubscriptionData
	}

	// Update fields
	subscription.Name = name
	subscription.StartDate = startDate
	subscription.DurationDays = durationDays
	subscription.NotificationEnabled = notificationEnabled

	if err := s.subscriptionRepo.Update(subscription); err != nil {
		return nil, err
	}

	return subscription, nil
}

func (s *subscriptionService) Delete(id, userID uint) error {
	// Verify ownership
	_, err := s.GetByID(id, userID)
	if err != nil {
		return err
	}

	return s.subscriptionRepo.Delete(id)
}

func (s *subscriptionService) ToggleNotification(id, userID uint, enabled bool) (*models.Subscription, error) {
	// Get existing subscription and verify ownership
	subscription, err := s.GetByID(id, userID)
	if err != nil {
		return nil, err
	}

	subscription.NotificationEnabled = enabled

	if err := s.subscriptionRepo.Update(subscription); err != nil {
		return nil, err
	}

	return subscription, nil
}
