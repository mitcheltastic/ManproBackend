package router

import (
	"context"
	"net/http"
	"strings"

	"firebase.google.com/go/v4/auth"
	"github.com/gin-gonic/gin"
	fbclient "github.com/mitcheltastic/ManproBackend/internal/infrastructure/firebase"
)

// ContextKey is used to store and retrieve data from context.
type ContextKey string

const AuthClaimsKey ContextKey = "firebaseAuthClaims"

func AuthMiddleware(fc *fbclient.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header required"})
			c.Abort()
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization format must be Bearer <token>"})
			c.Abort()
			return
		}

		idToken := parts[1]

		token, err := fc.AuthClient.VerifyIDToken(context.Background(), idToken)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired token", "details": err.Error()})
			c.Abort()
			return
		}

		c.Set(string(AuthClaimsKey), token)
		c.Set("userUID", token.UID)

		c.Next()
	}
}

func GetAuthClaims(c *gin.Context) (*auth.Token, bool) {
	value, exists := c.Get(string(AuthClaimsKey))
	if !exists {
		return nil, false
	}
	claims, ok := value.(*auth.Token)
	return claims, ok
}
