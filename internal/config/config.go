package config

import (
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	Server    ServerConfig
	Database  DatabaseConfig
	JWT       JWTConfig
	Email     EmailConfig
	Scheduler SchedulerConfig
}

type ServerConfig struct {
	Port   string
	Env    string
	GinMode string
}

type DatabaseConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
	SSLMode  string
}

type JWTConfig struct {
	Secret          string
	ExpirationHours int
}

type EmailConfig struct {
	SMTPHost     string
	SMTPPort     string
	SMTPUsername string
	SMTPPassword string
	FromEmail    string
	FromName     string
}

type SchedulerConfig struct {
	Enabled              bool
	CronExpression       string
	NotificationDaysBefore int
}

var AppConfig *Config

// Load reads configuration from environment variables
func Load() (*Config, error) {
	// Load .env file if it exists (ignore error in production)
	_ = godotenv.Load()

	// Validate required environment variables
	required := []string{
		"DB_HOST", "DB_PORT", "DB_USER", "DB_PASSWORD", "DB_NAME",
		"JWT_SECRET",
		"SMTP_HOST", "SMTP_PORT", "SMTP_USERNAME", "SMTP_PASSWORD",
	}

	for _, key := range required {
		if os.Getenv(key) == "" {
			return nil, fmt.Errorf("required environment variable %s is not set", key)
		}
	}

	jwtExpHours, err := strconv.Atoi(getEnv("JWT_EXPIRATION_HOURS", "72"))
	if err != nil {
		jwtExpHours = 72
	}

	schedulerEnabled, err := strconv.ParseBool(getEnv("SCHEDULER_ENABLED", "true"))
	if err != nil {
		schedulerEnabled = true
	}

	notificationDaysBefore, err := strconv.Atoi(getEnv("NOTIFICATION_DAYS_BEFORE", "5"))
	if err != nil {
		notificationDaysBefore = 5
	}

	config := &Config{
		Server: ServerConfig{
			Port:    getEnv("SERVER_PORT", "8080"),
			Env:     getEnv("APP_ENV", "development"),
			GinMode: getEnv("GIN_MODE", "debug"),
		},
		Database: DatabaseConfig{
			Host:     os.Getenv("DB_HOST"),
			Port:     os.Getenv("DB_PORT"),
			User:     os.Getenv("DB_USER"),
			Password: os.Getenv("DB_PASSWORD"),
			DBName:   os.Getenv("DB_NAME"),
			SSLMode:  getEnv("DB_SSLMODE", "disable"),
		},
		JWT: JWTConfig{
			Secret:          os.Getenv("JWT_SECRET"),
			ExpirationHours: jwtExpHours,
		},
		Email: EmailConfig{
			SMTPHost:     os.Getenv("SMTP_HOST"),
			SMTPPort:     os.Getenv("SMTP_PORT"),
			SMTPUsername: os.Getenv("SMTP_USERNAME"),
			SMTPPassword: os.Getenv("SMTP_PASSWORD"),
			FromEmail:    getEnv("SMTP_FROM_EMAIL", "noreply@renewguard.com"),
			FromName:     getEnv("SMTP_FROM_NAME", "RenewGuard"),
		},
		Scheduler: SchedulerConfig{
			Enabled:              schedulerEnabled,
			CronExpression:       getEnv("SCHEDULER_CRON", "0 0 * * *"),
			NotificationDaysBefore: notificationDaysBefore,
		},
	}

	AppConfig = config
	log.Println("Configuration loaded successfully")
	return config, nil
}

// GetDSN returns the PostgreSQL connection string
func (c *DatabaseConfig) GetDSN() string {
	return fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		c.Host, c.Port, c.User, c.Password, c.DBName, c.SSLMode,
	)
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
