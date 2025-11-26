package utils

import (
	"regexp"
	"strings"
)

// IsValidEmail validates email format
func IsValidEmail(email string) bool {
	email = strings.TrimSpace(email)
	if email == "" {
		return false
	}

	// Simple email regex
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`)
	return emailRegex.MatchString(email)
}

// IsValidPassword checks if password meets minimum requirements
func IsValidPassword(password string) bool {
	// Minimum 6 characters
	return len(password) >= 6
}
