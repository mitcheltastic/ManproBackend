package security

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

// Define custom claims struct that embeds the standard JWT claims
type UserClaims struct {
	UserID uuid.UUID `json:"user_id"`
	Email  string    `json:"email"`
	jwt.RegisteredClaims
}

// JWTService defines the interface for token operations.
type JWTService interface {
	GenerateToken(userID uuid.UUID, email string) (string, error)
}

// jwtServiceImpl is the concrete implementation of the JWTService.
type jwtServiceImpl struct {
	// IMPORTANT: In a real application, this should be loaded from a secret environment variable (e.g., config.JWTSecret).
	secretKey []byte 
	issuer    string
}

// NewJWTService creates a new JWT service instance.
func NewJWTService(secret string, issuer string) JWTService {
	return &jwtServiceImpl{
		secretKey: []byte(secret),
		issuer:    issuer,
	}
}

// GenerateToken creates a signed JWT string containing user information.
func (s *jwtServiceImpl) GenerateToken(userID uuid.UUID, email string) (string, error) {
	// Define the token claims
	claims := UserClaims{
		UserID: userID,
		Email:  email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 24)), // Token expires in 24 hours
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Issuer:    s.issuer,
			Subject:   userID.String(),
		},
	}

	// Create the token using the claims and the HS256 signing method
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Sign the token with the secret key
	tokenString, err := token.SignedString(s.secretKey)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}