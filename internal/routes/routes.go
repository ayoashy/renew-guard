package routes

import (
	"net/http"
	"renew-guard/internal/controllers"
	"renew-guard/internal/middleware"
	"renew-guard/pkg/jwt"

	"github.com/gin-gonic/gin"
)

type Router struct {
	authController         *controllers.AuthController
	subscriptionController *controllers.SubscriptionController
	emailTestController    *controllers.EmailTestController
	jwtUtil                *jwt.JWTUtil
}

func NewRouter(
	authController *controllers.AuthController,
	subscriptionController *controllers.SubscriptionController,
	emailTestController *controllers.EmailTestController,
	jwtUtil *jwt.JWTUtil,
) *Router {
	return &Router{
		authController:         authController,
		subscriptionController: subscriptionController,
		emailTestController:    emailTestController,
		jwtUtil:                jwtUtil,
	}
}

// SetupRoutes configures all application routes
func (r *Router) SetupRoutes(router *gin.Engine) {
	// Apply global middleware
	router.Use(middleware.CORSMiddleware())
	router.Use(middleware.ErrorMiddleware())

	// Health check endpoint
	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status": "healthy",
			"service": "renew-guard",
		})
	})

	// API v1 routes
	api := router.Group("/api/v1")
	{
		// Auth routes (public)
		auth := api.Group("/auth")
		{
			auth.POST("/register", r.authController.Register)
			auth.POST("/login", r.authController.Login)
		}

		// Subscription routes (protected)
		subscriptions := api.Group("/subscriptions")
		subscriptions.Use(middleware.AuthMiddleware(r.jwtUtil))
		{
			subscriptions.POST("", r.subscriptionController.CreateSubscription)
			subscriptions.GET("", r.subscriptionController.GetSubscriptions)
			subscriptions.GET("/:id", r.subscriptionController.GetSubscription)
			subscriptions.PUT("/:id", r.subscriptionController.UpdateSubscription)
			subscriptions.DELETE("/:id", r.subscriptionController.DeleteSubscription)
			subscriptions.PATCH("/:id/notifications", r.subscriptionController.ToggleNotification)
		}

		// Email test routes (public - for testing SMTP)
		test := api.Group("/test")
		{
			test.POST("/email", r.emailTestController.SendTestEmail)
		}
	}
}
