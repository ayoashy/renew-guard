package controllers

import (
	"net/http"
	"renew-guard/internal/services"
	"renew-guard/pkg/utils"

	"github.com/gin-gonic/gin"
)

type AuthController struct {
	authService services.AuthService
}

func NewAuthController(authService services.AuthService) *AuthController {
	return &AuthController{
		authService: authService,
	}
}

type RegisterRequest struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type LoginRequest struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type AuthResponse struct {
	User  UserResponse `json:"user"`
	Token string       `json:"token"`
}

type UserResponse struct {
	ID    uint   `json:"id"`
	Email string `json:"email"`
}

// Register handles user registration
// @Summary Register a new user
// @Tags auth
// @Accept json
// @Produce json
// @Param request body RegisterRequest true "Registration details"
// @Success 201 {object} AuthResponse
// @Router /api/auth/register [post]
func (ctrl *AuthController) Register(c *gin.Context) {
	var req RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid request body")
		return
	}

	user, token, err := ctrl.authService.Register(req.Email, req.Password)
	if err != nil {
		switch err {
		case services.ErrEmailAlreadyExists:
			utils.ErrorResponse(c, http.StatusConflict, "Email already exists")
		case services.ErrInvalidEmail:
			utils.ErrorResponse(c, http.StatusBadRequest, "Invalid email format")
		case services.ErrWeakPassword:
			utils.ErrorResponse(c, http.StatusBadRequest, "Password must be at least 6 characters")
		default:
			utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to register user")
		}
		return
	}

	response := AuthResponse{
		User: UserResponse{
			ID:    user.ID,
			Email: user.Email,
		},
		Token: token,
	}

	utils.SuccessResponse(c, http.StatusCreated, "User registered successfully", response)
}

// Login handles user authentication
// @Summary Login user
// @Tags auth
// @Accept json
// @Produce json
// @Param request body LoginRequest true "Login credentials"
// @Success 200 {object} AuthResponse
// @Router /api/auth/login [post]
func (ctrl *AuthController) Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid request body")
		return
	}

	user, token, err := ctrl.authService.Login(req.Email, req.Password)
	if err != nil {
		if err == services.ErrInvalidCredentials {
			utils.ErrorResponse(c, http.StatusUnauthorized, "Invalid email or password")
		} else {
			utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to login")
		}
		return
	}

	response := AuthResponse{
		User: UserResponse{
			ID:    user.ID,
			Email: user.Email,
		},
		Token: token,
	}

	utils.SuccessResponse(c, http.StatusOK, "Login successful", response)
}
