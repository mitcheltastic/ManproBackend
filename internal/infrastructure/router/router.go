package router

import (
	"net/http"

	"github.com/gin-gonic/gin"
	dbimpl "github.com/mitcheltastic/ManproBackend/internal/infrastructure/database"
	fbclient "github.com/mitcheltastic/ManproBackend/internal/infrastructure/firebase" 
)

// SetupRoutes registers all API routes and middleware.
// We now accept the database client to pass to handlers/controllers.
func SetupRoutes(r *gin.Engine, fbClient *fbclient.Client, dbClient *dbimpl.Client) {
	
	// --- Public Routes ---
	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "OK", "service": "Go Backend"})
	})

	// --- V1 API Group ---
	v1 := r.Group("/api/v1")
	
	// Future: Auth Controller initialization will go here, using dbClient:
	// authController := handler.NewAuthController(dbClient)

	// Public Auth Endpoints (No token required)
	// v1.POST("/auth/register", authController.Register)
	// v1.POST("/auth/login", authController.Login)
	// v1.POST("/auth/forgot-password", authController.ForgotPassword)

	// Protected Routes (Require Firebase Authentication Middleware)
	v1.Use(AuthMiddleware(fbClient))
	{
		// Example protected endpoint: Get User Profile
		// This handler will only run if AuthMiddleware successfully verified the Firebase ID token
		v1.GET("/profile", protectedProfileHandler)
		
		// Future: Authenticated endpoints like CRUD APIs will go here
	}
}

// protectedProfileHandler is a sample handler that retrieves the user's information
// from the request context after the token has been verified.
func protectedProfileHandler(c *gin.Context) {
	// Retrieve the decoded claims injected by the AuthMiddleware
	claims, ok := GetAuthClaims(c)
	if !ok {
		// Should not happen if middleware ran, but defensive programming
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Authentication context missing"})
		return
	}

	// Return the user data retrieved from the Firebase token
	c.JSON(http.StatusOK, gin.H{
		"message": "Welcome! You are authenticated.",
		"user_id": claims.UID,
		"email":   claims.Claims["email"], 
		"name":    claims.Claims["name"], 
		"claims_raw": claims.Claims,
	})
}