package router

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mitcheltastic/ManproBackend/config" // Import for config
	"github.com/mitcheltastic/ManproBackend/internal/handler" 
	"github.com/mitcheltastic/ManproBackend/internal/pkg/security" 
	"github.com/mitcheltastic/ManproBackend/internal/pkg/email"     // New Import for Email Sender
	"github.com/mitcheltastic/ManproBackend/internal/service" 
	dbimpl "github.com/mitcheltastic/ManproBackend/internal/infrastructure/database"
	fbclient "github.com/mitcheltastic/ManproBackend/internal/infrastructure/firebase" 
)

// SetupRoutes registers all API routes and middleware.
func SetupRoutes(r *gin.Engine, fbClient *fbclient.Client, dbClient *dbimpl.Client, cfg *config.Config) {
	
	// --- Dependency Injection Setup (Wiring the Layers) ---
	
	// 1. Initialize Repository (Data Access)
	authRepo := dbimpl.NewAuthRepository(dbClient.DB) 

	// 2. Initialize JWT Service (The Token Generator)
	jwtService := security.NewJWTService(cfg.JWTSecret, "manpro_backend")

	// 3. Initialize Email Sender (The critical new piece)
	emailSender := email.NewSMTPSender(
		cfg.SMTPHost,
		cfg.SMTPPort,
		cfg.SMTPUser,
		cfg.SMTPPass,
		cfg.FromEmail,
	)

	// 4. Initialize Service (Business Logic)
	// FIX: Pass the emailSender as the third argument
	authService := service.NewAuthService(authRepo, jwtService, emailSender) 

	// 5. Initialize Handler (HTTP Controller)
	authHandler := handler.NewAuthHandler(authService)

	// --- Public Routes ---
	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "OK", "service": "Go Backend"})
	})

	// --- V1 API Group ---
	v1 := r.Group("/api/v1")
	
	// Public Auth Endpoints
	v1.POST("/auth/register", authHandler.Register)
	v1.POST("/auth/login", authHandler.Login)
	v1.POST("/auth/forgot-password", authHandler.StartPasswordReset)
	v1.POST("/auth/reset-password", authHandler.ResetPassword)

	// Protected Routes (Require Firebase Authentication Middleware)
	v1.Use(AuthMiddleware(fbClient))
	{
		// Example protected endpoint: Get User Profile (using Firebase claims)
		v1.GET("/profile", protectedProfileHandler)
	}
}

// protectedProfileHandler is a sample handler that retrieves the user's information
// from the request context after the token has been verified.
func protectedProfileHandler(c *gin.Context) {
	claims, ok := GetAuthClaims(c)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Authentication context missing"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Welcome! You are authenticated.",
		"user_id": claims.UID,
		"email":   claims.Claims["email"], 
		"name":    claims.Claims["name"], 
		"claims_raw": claims.Claims,
	})
}