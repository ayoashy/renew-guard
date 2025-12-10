package controllers

import (
	"net/http"
	"renew-guard/internal/middleware"
	"renew-guard/internal/services"
	"renew-guard/pkg/email"
	"renew-guard/pkg/utils"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

type SubscriptionController struct {
	subscriptionService services.SubscriptionService
	emailService        email.EmailService
}

func NewSubscriptionController(subscriptionService services.SubscriptionService, emailService email.EmailService) *SubscriptionController {
	return &SubscriptionController{
		subscriptionService: subscriptionService,
		emailService:        emailService,
	}
}

type CreateSubscriptionRequest struct {
	Name         string    `json:"name" binding:"required"`
	StartDate    time.Time `json:"start_date" binding:"required"`
	DurationDays int       `json:"duration_days" binding:"required,min=1"`
}

type UpdateSubscriptionRequest struct {
	Name                string    `json:"name" binding:"required"`
	StartDate           time.Time `json:"start_date" binding:"required"`
	DurationDays        int       `json:"duration_days" binding:"required,min=1"`
	NotificationEnabled bool      `json:"notification_enabled"`
}

type ToggleNotificationRequest struct {
	Enabled bool `json:"enabled"`
}

// CreateSubscription creates a new subscription
// @Summary Create a new subscription
// @Tags subscriptions
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body CreateSubscriptionRequest true "Subscription details"
// @Success 201 {object} models.Subscription
// @Router /api/subscriptions [post]
func (ctrl *SubscriptionController) CreateSubscription(c *gin.Context) {
	userID, exists := middleware.GetUserID(c)
	if !exists {
		utils.ErrorResponse(c, http.StatusUnauthorized, "Unauthorized")
		return
	}

	// Get user email from context
	userEmail, exists := middleware.GetUserEmail(c)
	if !exists {
		utils.ErrorResponse(c, http.StatusUnauthorized, "User email not found")
		return
	}

	var req CreateSubscriptionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid request body")
		return
	}

	// Create subscription with user's email
	subscription, err := ctrl.subscriptionService.Create(userID, userEmail, req.Name, req.StartDate, req.DurationDays)
	if err != nil {
		if err == services.ErrInvalidSubscriptionData {
			utils.ErrorResponse(c, http.StatusBadRequest, "Invalid subscription data")
		} else {
			utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to create subscription")
		}
		return
	}

	// Send confirmation email asynchronously
	go ctrl.sendSubscriptionConfirmation(userEmail, subscription.Name, subscription.StartDate, subscription.EndDate)

	utils.SuccessResponse(c, http.StatusCreated, "Subscription created successfully", subscription)
}

// GetSubscriptions retrieves all subscriptions for the authenticated user
// @Summary Get all user subscriptions
// @Tags subscriptions
// @Produce json
// @Security BearerAuth
// @Success 200 {array} models.Subscription
// @Router /api/subscriptions [get]
func (ctrl *SubscriptionController) GetSubscriptions(c *gin.Context) {
	userID, exists := middleware.GetUserID(c)
	if !exists {
		utils.ErrorResponse(c, http.StatusUnauthorized, "Unauthorized")
		return
	}

	subscriptions, err := ctrl.subscriptionService.GetAllByUserID(userID)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to retrieve subscriptions")
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Subscriptions retrieved successfully", subscriptions)
}

// GetSubscription retrieves a specific subscription by ID
// @Summary Get subscription by ID
// @Tags subscriptions
// @Produce json
// @Security BearerAuth
// @Param id path int true "Subscription ID"
// @Success 200 {object} models.Subscription
// @Router /api/subscriptions/{id} [get]
func (ctrl *SubscriptionController) GetSubscription(c *gin.Context) {
	userID, exists := middleware.GetUserID(c)
	if !exists {
		utils.ErrorResponse(c, http.StatusUnauthorized, "Unauthorized")
		return
	}

	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid subscription ID")
		return
	}

	subscription, err := ctrl.subscriptionService.GetByID(uint(id), userID)
	if err != nil {
		if err == services.ErrSubscriptionNotFound {
			utils.ErrorResponse(c, http.StatusNotFound, "Subscription not found")
		} else if err == services.ErrUnauthorizedAccess {
			utils.ErrorResponse(c, http.StatusForbidden, "Unauthorized access")
		} else {
			utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to retrieve subscription")
		}
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Subscription retrieved successfully", subscription)
}

// UpdateSubscription updates an existing subscription
// @Summary Update subscription
// @Tags subscriptions
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Subscription ID"
// @Param request body UpdateSubscriptionRequest true "Updated subscription details"
// @Success 200 {object} models.Subscription
// @Router /api/subscriptions/{id} [put]
func (ctrl *SubscriptionController) UpdateSubscription(c *gin.Context) {
	userID, exists := middleware.GetUserID(c)
	if !exists {
		utils.ErrorResponse(c, http.StatusUnauthorized, "Unauthorized")
		return
	}

	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid subscription ID")
		return
	}

	var req UpdateSubscriptionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid request body")
		return
	}

	subscription, err := ctrl.subscriptionService.Update(
		uint(id), userID, req.Name, req.StartDate, req.DurationDays, req.NotificationEnabled,
	)
	if err != nil {
		if err == services.ErrSubscriptionNotFound {
			utils.ErrorResponse(c, http.StatusNotFound, "Subscription not found")
		} else if err == services.ErrUnauthorizedAccess {
			utils.ErrorResponse(c, http.StatusForbidden, "Unauthorized access")
		} else if err == services.ErrInvalidSubscriptionData {
			utils.ErrorResponse(c, http.StatusBadRequest, "Invalid subscription data")
		} else {
			utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to update subscription")
		}
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Subscription updated successfully", subscription)
}

// DeleteSubscription deletes a subscription
// @Summary Delete subscription
// @Tags subscriptions
// @Produce json
// @Security BearerAuth
// @Param id path int true "Subscription ID"
// @Success 200 {object} Response
// @Router /api/subscriptions/{id} [delete]
func (ctrl *SubscriptionController) DeleteSubscription(c *gin.Context) {
	userID, exists := middleware.GetUserID(c)
	if !exists {
		utils.ErrorResponse(c, http.StatusUnauthorized, "Unauthorized")
		return
	}

	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid subscription ID")
		return
	}

	err = ctrl.subscriptionService.Delete(uint(id), userID)
	if err != nil {
		if err == services.ErrSubscriptionNotFound {
			utils.ErrorResponse(c, http.StatusNotFound, "Subscription not found")
		} else if err == services.ErrUnauthorizedAccess {
			utils.ErrorResponse(c, http.StatusForbidden, "Unauthorized access")
		} else {
			utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to delete subscription")
		}
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Subscription deleted successfully", nil)
}

// ToggleNotification toggles notification settings for a subscription
// @Summary Toggle subscription notifications
// @Tags subscriptions
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Subscription ID"
// @Param request body ToggleNotificationRequest true "Notification settings"
// @Success 200 {object} models.Subscription
// @Router /api/subscriptions/{id}/notifications [patch]
func (ctrl *SubscriptionController) ToggleNotification(c *gin.Context) {
	userID, exists := middleware.GetUserID(c)
	if !exists {
		utils.ErrorResponse(c, http.StatusUnauthorized, "Unauthorized")
		return
	}

	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid subscription ID")
		return
	}

	var req ToggleNotificationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid request body")
		return
	}

	subscription, err := ctrl.subscriptionService.ToggleNotification(uint(id), userID, req.Enabled)
	if err != nil {
		if err == services.ErrSubscriptionNotFound {
			utils.ErrorResponse(c, http.StatusNotFound, "Subscription not found")
		} else if err == services.ErrUnauthorizedAccess {
			utils.ErrorResponse(c, http.StatusForbidden, "Unauthorized access")
		} else {
			utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to toggle notification")
		}
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Notification settings updated successfully", subscription)
}
