package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"renew-guard/internal/config"
	"renew-guard/internal/controllers"
	"renew-guard/internal/database"
	"renew-guard/internal/repositories"
	"renew-guard/internal/routes"
	"renew-guard/internal/scheduler"
	"renew-guard/internal/services"
	"renew-guard/pkg/email"
	"renew-guard/pkg/jwt"
	"syscall"

	"github.com/gin-gonic/gin"
)

func main() {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Set Gin mode
	gin.SetMode(cfg.Server.GinMode)

	// Initialize database
	db, err := database.Initialize(cfg)
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer database.Close()

	// Initialize repositories
	userRepo := repositories.NewUserRepository(db)
	subscriptionRepo := repositories.NewSubscriptionRepository(db)
	notificationLogRepo := repositories.NewNotificationLogRepository(db)

	// Initialize JWT utility
	jwtUtil := jwt.NewJWTUtil(cfg.JWT.Secret, cfg.JWT.ExpirationHours)

	// Initialize email service
	emailConfig := email.EmailConfig{
		SMTPHost:     cfg.Email.SMTPHost,
		SMTPPort:     cfg.Email.SMTPPort,
		SMTPUsername: cfg.Email.SMTPUsername,
		SMTPPassword: cfg.Email.SMTPPassword,
		FromEmail:    cfg.Email.FromEmail,
		FromName:     cfg.Email.FromName,
	}
	emailService := email.NewEmailService(emailConfig)

	// Initialize services
	authService := services.NewAuthService(userRepo, jwtUtil)
	subscriptionService := services.NewSubscriptionService(subscriptionRepo)
	notificationService := services.NewNotificationService(subscriptionRepo, notificationLogRepo, emailService)

	// Initialize controllers
	authController := controllers.NewAuthController(authService)
	subscriptionController := controllers.NewSubscriptionController(subscriptionService)
	emailTestController := controllers.NewEmailTestController(emailService)

	// Initialize router
	router := gin.Default()
	appRouter := routes.NewRouter(authController, subscriptionController, emailTestController, jwtUtil)
	appRouter.SetupRoutes(router)

	// Initialize and start scheduler
	schedulerInstance := scheduler.NewScheduler(notificationService, &cfg.Scheduler)
	if err := schedulerInstance.Start(); err != nil {
		log.Fatalf("Failed to start scheduler: %v", err)
	}
	defer schedulerInstance.Stop()

	// Start HTTP server
	serverAddr := fmt.Sprintf(":%s", cfg.Server.Port)
	log.Printf("Starting server on %s", serverAddr)
	log.Printf("Environment: %s", cfg.Server.Env)

	// Graceful shutdown
	go func() {
		if err := router.Run(serverAddr); err != nil {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	// Wait for interrupt signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down server...")
}
