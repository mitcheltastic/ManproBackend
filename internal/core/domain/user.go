package domain

import (
	"time"

	"github.com/google/uuid"
)

// User represents the core user model stored in the database (Supabase/Postgres).
type User struct {
	// --- CORE AUTHENTICATION FIELDS (Required for Login/Security) ---
	ID        uuid.UUID `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	HashedPassword string `json:"-"`
	IsVerified    bool   `json:"is_verified"`
	
	// --- PROFILE FIELDS (The new Intertwine data) ---
	Nickname      *string    `json:"nickname"`      // Optional
	Gender        *string    `json:"gender"`        // M, F, Other
	Birthdate     *time.Time `json:"birthdate"`     // Used for age calculation
	College       *string    `json:"college"`       // e.g., "MIT"
	Faculty       *string    `json:"faculty"`       // e.g., "Engineering"
	Major         *string    `json:"major"`         // e.g., "Computer Science"
	Year          *int       `json:"year"`          // Current academic year (e.g., 2024)
	MBTI          *string    `json:"mbti"`          // Personality type
	BloodType     *string    `json:"blood_type"`    // A, B, AB, O
	ProfilePictureURL *string `json:"profile_picture_url"` // S3 or Storage URL for profile image
	
	// GalleryURLs stores up to 5 image URLs. We will use a string slice in Go.
	// NOTE: In Postgres/Supabase, this will be stored as TEXT[] (an array of text/strings).
	GalleryPictureURLs []string `json:"gallery_picture_urls"`

	// Hobby stores up to 3 hobbies as a slice.
	Hobby []string `json:"hobby"`

	// --- METADATA FIELDS ---
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// --- Request/Input Models (DTOs) ---

// RegisterRequest holds the user input for the registration endpoint.
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