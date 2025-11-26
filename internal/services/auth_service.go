package services

import (
	"errors"
	"renew-guard/internal/models"
	"renew-guard/internal/repositories"
	"renew-guard/pkg/jwt"
	"renew-guard/pkg/utils"

	"gorm.io/gorm"
)

var (
	ErrEmailAlreadyExists = errors.New("email already exists")
	ErrInvalidCredentials = errors.New("invalid email or password")
	ErrInvalidEmail       = errors.New("invalid email format")
	ErrWeakPassword       = errors.New("password must be at least 6 characters")
)

type AuthService interface {
	Register(email, password string) (*models.User, string, error)
	Login(email, password string) (*models.User, string, error)
}

type authService struct {
	userRepo repositories.UserRepository
	jwtUtil  *jwt.JWTUtil
}

func NewAuthService(userRepo repositories.UserRepository, jwtUtil *jwt.JWTUtil) AuthService {
	return &authService{
		userRepo: userRepo,
		jwtUtil:  jwtUtil,
	}
}

func (s *authService) Register(email, password string) (*models.User, string, error) {
	// Validate email
	if !utils.IsValidEmail(email) {
		return nil, "", ErrInvalidEmail
	}

	// Validate password
	if !utils.IsValidPassword(password) {
		return nil, "", ErrWeakPassword
	}

	// Check if user already exists
	existingUser, err := s.userRepo.FindByEmail(email)
	if err == nil && existingUser != nil {
		return nil, "", ErrEmailAlreadyExists
	}

	// Create new user
	user := &models.User{
		Email: email,
	}

	if err := user.HashPassword(password); err != nil {
		return nil, "", err
	}

	if err := s.userRepo.Create(user); err != nil {
		return nil, "", err
	}

	// Generate JWT token
	token, err := s.jwtUtil.GenerateToken(user.ID, user.Email)
	if err != nil {
		return nil, "", err
	}

	return user, token, nil
}

func (s *authService) Login(email, password string) (*models.User, string, error) {
	// Find user by email
	user, err := s.userRepo.FindByEmail(email)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, "", ErrInvalidCredentials
		}
		return nil, "", err
	}

	// Check password
	if !user.CheckPassword(password) {
		return nil, "", ErrInvalidCredentials
	}

	// Generate JWT token
	token, err := s.jwtUtil.GenerateToken(user.ID, user.Email)
	if err != nil {
		return nil, "", err
	}

	return user, token, nil
}
