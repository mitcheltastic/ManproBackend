package domain

import (
	"time"

	"github.com/google/uuid"
)

// User represents the core user model stored in the database (Supabase/Postgres).
// This struct will form the basis of your 'users' table in Supabase.
type User struct {
	ID        uuid.UUID `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	// HashedPassword stores the bcrypt hash. We never expose the raw password.
	HashedPassword string `json:"-"`
	IsVerified    bool   `json:"is_verified"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// --- Request/Input Models (DTOs) ---

// RegisterRequest holds the user input for the registration endpoint.
// The 'binding' tags are used by Gin/validator to enforce input rules.
type RegisterRequest struct {
	Name            string `json:"name" binding:"required,min=2,max=50"`
	Email           string `json:"email" binding:"required,email"`
	Password        string `json:"password" binding:"required,min=8"`
	ConfirmPassword string `json:"confirm_password" binding:"required,eqfield=Password"`
}

// LoginRequest holds the user input for the login endpoint.
type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

// ForgotPasswordRequest holds the email input for initiating the forgot password flow.
type ForgotPasswordRequest struct {
	Email string `json:"email" binding:"required,email"`
}

// ResetPasswordRequest holds the input for completing the password reset.
type ResetPasswordRequest struct {
	Email           string `json:"email" binding:"required,email"`
	Code            string `json:"code" binding:"required,len=6"` // Assuming a 6-digit code
	NewPassword     string `json:"new_password" binding:"required,min=8"`
	ConfirmPassword string `json:"confirm_password" binding:"required,eqfield=NewPassword"`
}

// --- Response Models ---

// AuthResponse holds the data returned to the client upon successful registration or login.
type AuthResponse struct {
	UserID uuid.UUID `json:"user_id"`
	Email  string    `json:"email"`
	Name   string    `json:"name"`
	Token  string    `json:"token"` // The JWT issued to the client
}