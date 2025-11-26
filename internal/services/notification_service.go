package services

import (
	"fmt"
	"log"
	"renew-guard/internal/models"
	"renew-guard/internal/repositories"
	"renew-guard/pkg/email"
	"time"
)

type NotificationService interface {
	CheckAndSendNotifications(daysBefore int) error
	SendExpirationWarning(subscription *models.Subscription) error
}

type notificationService struct {
	subscriptionRepo repositories.SubscriptionRepository
	notificationRepo repositories.NotificationLogRepository
	emailService     email.EmailService
}

func NewNotificationService(
	subscriptionRepo repositories.SubscriptionRepository,
	notificationRepo repositories.NotificationLogRepository,
	emailService email.EmailService,
) NotificationService {
	return &notificationService{
		subscriptionRepo: subscriptionRepo,
		notificationRepo: notificationRepo,
		emailService:     emailService,
	}
}

func (s *notificationService) CheckAndSendNotifications(daysBefore int) error {
	log.Printf("Checking for subscriptions expiring within %d days...", daysBefore)

	// Find all subscriptions that need notification
	subscriptions, err := s.subscriptionRepo.FindExpiringSubscriptions(daysBefore)
	if err != nil {
		return fmt.Errorf("failed to find expiring subscriptions: %w", err)
	}

	log.Printf("Found %d subscription(s) requiring notification", len(subscriptions))

	sentCount := 0
	failedCount := 0

	for _, subscription := range subscriptions {
		// Check if should notify (includes daily check)
		if subscription.ShouldNotify(daysBefore) {
			if err := s.SendExpirationWarning(&subscription); err != nil {
				log.Printf("Failed to send notification for subscription %d: %v", subscription.ID, err)
				failedCount++
			} else {
				sentCount++
			}
		}
	}

	log.Printf("Notification run complete: %d sent, %d failed", sentCount, failedCount)
	return nil
}

func (s *notificationService) SendExpirationWarning(subscription *models.Subscription) error {
	daysLeft := subscription.DaysUntilExpiration()

	// Generate email content
	subject := email.GetExpirationWarningSubject(subscription.Name, daysLeft)
	htmlBody := email.GetExpirationWarningTemplate(subscription.Name, daysLeft, subscription.EndDate)

	// Send email
	err := s.emailService.SendHTML(subscription.User.Email, subject, htmlBody)

	// Create notification log
	notificationLog := &models.NotificationLog{
		SubscriptionID: subscription.ID,
		Status:         "success",
	}

	if err != nil {
		notificationLog.Status = "failed"
		notificationLog.ErrorMessage = err.Error()
		
		// Still log the failed attempt
		if logErr := s.notificationRepo.Create(notificationLog); logErr != nil {
			log.Printf("Failed to create notification log: %v", logErr)
		}
		
		return fmt.Errorf("failed to send email to %s: %w", subscription.User.Email, err)
	}

	// Update last notification sent timestamp
	now := time.Now()
	if err := s.subscriptionRepo.UpdateLastNotificationSent(subscription.ID, now); err != nil {
		log.Printf("Failed to update last notification sent: %v", err)
	}

	// Create success log
	if err := s.notificationRepo.Create(notificationLog); err != nil {
		log.Printf("Failed to create notification log: %v", err)
	}

	log.Printf("Notification sent successfully for subscription %d (%s) to %s", 
		subscription.ID, subscription.Name, subscription.User.Email)

	return nil
}
