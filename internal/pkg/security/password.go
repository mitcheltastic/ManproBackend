package security

import (
	"crypto/rand"
	"log"
	"math/big"

	"golang.org/x/crypto/bcrypt"
)

// HashPassword takes a plaintext password and returns its bcrypt hash.
func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}

// CheckPasswordHash compares a plaintext password with a hashed password.
// Returns nil on success, or an error if they do not match.
func CheckPasswordHash(password, hash string) error {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
}

// GenerateNumericCode generates a cryptographically secure, random numeric string of the given length.
func GenerateNumericCode(length int) string {
	const charset = "0123456789"
	b := make([]byte, length)
	
	for i := range b {
		// Generate a random index within the charset length
		randomIndex, err := rand.Int(rand.Reader, big.NewInt(int64(len(charset))))
		if err != nil {
			// In a real application, this should be handled more gracefully, but for simplicity:
			log.Fatalf("Failed to generate random number: %v", err)
		}
		b[i] = charset[randomIndex.Int64()]
	}
	return string(b)
}