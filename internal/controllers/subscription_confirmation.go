package controllers

import (
	"log"
	"renew-guard/pkg/email"
	"time"
)

// sendSubscriptionConfirmation sends a confirmation email when a subscription is created
func (ctrl *SubscriptionController) sendSubscriptionConfirmation(userEmail, subscriptionName string, startDate, endDate time.Time) {
	subject := email.GetSubscriptionConfirmationSubject(subscriptionName)
	htmlBody := email.GetSubscriptionConfirmationTemplate(subscriptionName, startDate, endDate)
	
	err := ctrl.emailService.SendHTML(userEmail, subject, htmlBody)
	if err != nil {
		log.Printf("Failed to send subscription confirmation email to %s: %v", userEmail, err)
	} else {
		log.Printf("Subscription confirmation email sent to %s for %s", userEmail, subscriptionName)
	}
}
