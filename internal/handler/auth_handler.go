package handler

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mitcheltastic/ManproBackend/internal/core/domain"
	"github.com/mitcheltastic/ManproBackend/internal/core/ports"
)

// AuthHandler handles HTTP requests related to standard authentication (Register, Login).
type AuthHandler struct {
	AuthService ports.AuthService
}

// NewAuthHandler creates a new instance of the AuthHandler.
func NewAuthHandler(authService ports.AuthService) *AuthHandler {
	return &AuthHandler{
		AuthService: authService,
	}
}

// Register handles user registration request (POST /api/v1/auth/register)
func (h *AuthHandler) Register(c *gin.Context) {
	var req domain.RegisterRequest
	
	// Bind JSON request body to the RegisterRequest struct and validate fields
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input", "details": err.Error()})
		return
	}

	// Password confirmation is handled by the 'eqfield' binding tag in the domain model,
	// but we can add an explicit check here if the binding somehow failed silently.
	if req.Password != req.ConfirmPassword {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Password and confirmation do not match"})
		return
	}

	// Call the service layer business logic
	authResponse, err := h.AuthService.Register(c, req)
	if err != nil {
		log.Printf("Registration error: %v", err)
		// Check for specific errors (e.g., duplicate email)
		if err.Error() == "a user with this email already exists" {
			c.JSON(http.StatusConflict, gin.H{"error": "User registration failed", "details": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "User registration failed"})
		return
	}

	c.JSON(http.StatusCreated, authResponse)
}

// Login handles user login request (POST /api/v1/auth/login)
func (h *AuthHandler) Login(c *gin.Context) {
	var req domain.LoginRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input", "details": err.Error()})
		return
	}

	authResponse, err := h.AuthService.Login(c, req)
	if err != nil {
		// Use generic message for security if credentials fail
		if err.Error() == "invalid credentials" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email or password"})
			return
		}
		log.Printf("Login error: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Login failed"})
		return
	}

	c.JSON(http.StatusOK, authResponse)
}

// StartPasswordReset handles initiating the forgot password flow (POST /api/v1/auth/forgot-password)
// REVISED: Now correctly captures the 'code' and 'err' return values.
func (h *AuthHandler) StartPasswordReset(c *gin.Context) {
	var req domain.ForgotPasswordRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input", "details": err.Error()})
		return
	}

	// FIX: Correctly capture both the code and the error from the service layer call
	code, err := h.AuthService.StartPasswordReset(c, req)
	if err != nil {
		log.Printf("Password reset initiation error: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to initiate password reset"})
		return
	}
	
	// IMPORTANT: Always return success (200 OK) even if the user is not found (code will be empty).
	// We return the code only for local testing.
	response := gin.H{"message": "If the email is registered, a password reset code has been sent."}
	if code != "" {
		response["debug_code"] = code // Include code for testing purposes
	}

	c.JSON(http.StatusOK, response)
}

// ResetPassword handles completing the password reset using the code (POST /api/v1/auth/reset-password)
func (h *AuthHandler) ResetPassword(c *gin.Context) {
	var req domain.ResetPasswordRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input", "details": err.Error()})
		return
	}

	authResponse, err := h.AuthService.ResetPassword(c, req)
	if err != nil {
		// Return specific error messages for user feedback on reset failure
		log.Printf("Password reset error: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Password reset failed", "details": err.Error()})
		return
	}

	c.JSON(http.StatusOK, authResponse)
}