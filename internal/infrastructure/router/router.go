package router

import (
	"net/http"

	"github.com/gin-gonic/gin"
	dbimpl "github.com/mitcheltastic/ManproBackend/internal/infrastructure/database" // Import for the DB Client struct
	fbclient "github.com/mitcheltastic/ManproBackend/internal/infrastructure/firebase" 
)

// SetupRoutes registers all API routes and middleware.
// NOTE: We now accept the dbClient to pass to handlers/controllers.
func SetupRoutes(r *gin.Engine, fbClient *fbclient.Client, dbClient *dbimpl.Client) {
	// Public routes (No authentication required)
	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "OK", "service": "Go Backend"})
	})

	// API Group V1
	v1 := r.Group("/api/v1")
	
	// Protected routes (Require Firebase Authentication Middleware)
	// Apply the middleware to all routes in this group
	v1.Use(AuthMiddleware(fbClient))
	{
		// Example protected endpoint: Get User Profile
		// This handler will only run if AuthMiddleware successfully verified the Firebase ID token
		v1.GET("/profile", protectedProfileHandler)
		
		// Future CRUD APIs will go here...
		// v1.POST("/users", userHandler.Create)
	}
}

// protectedProfileHandler is a sample handler that retrieves the user's information
// from the request context after the token has been verified.
func protectedProfileHandler(c *gin.Context) {
	// Retrieve the decoded claims injected by the AuthMiddleware
	claims, ok := GetAuthClaims(c)
	if !ok {
		// Should not happen if middleware ran, but good defensive programming
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Authentication context missing"})
		return
	}

	// Use the claims to return user data
	c.JSON(http.StatusOK, gin.H{
		"message": "Welcome! You are authenticated.",
		"user_id": claims.UID,
		"email":   claims.Claims["email"], // Accessing email claim directly
		"name":    claims.Claims["name"],   // Accessing name claim directly
		"claims_raw": claims.Claims,
	})
}