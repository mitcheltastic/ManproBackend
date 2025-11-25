package ports

import (
	"context"

	"github.com/mitcheltastic/ManproBackend/internal/core/domain"
)

// AuthRepository defines the interface for data operations related to authentication.
// The implementation will live in internal/infrastructure/database/
type AuthRepository interface {
	// CreateUser saves a new user record to the database (Supabase).
	CreateUser(ctx context.Context, user domain.User) error

	// GetUserByEmail retrieves a user by their email address.
	GetUserByEmail(ctx context.Context, email string) (*domain.User, error)

	// UpdateUserPassword updates the user's password hash in the database.
	UpdateUserPassword(ctx context.Context, userID string, newHashedPassword string) error
	
	// CreatePasswordResetCode saves a unique 6-digit code linked to a user/email.
	CreatePasswordResetCode(ctx context.Context, email string, code string) error
	
	// VerifyPasswordResetCode checks if the code is valid and not expired.
	VerifyPasswordResetCode(ctx context.Context, email string, code string) error
	
	// DeletePasswordResetCode removes the code after a successful reset.
	DeletePasswordResetCode(ctx context.Context, email string) error
}

// AuthService defines the interface for core business logic related to authentication.
// The implementation will live in internal/service/
type AuthService interface {
	// Register handles validation, hashing, saving the user, and issuing a token.
	Register(ctx context.Context, req domain.RegisterRequest) (*domain.AuthResponse, error)

	// Login verifies credentials and issues a token.
	Login(ctx context.Context, req domain.LoginRequest) (*domain.AuthResponse, error)

	// StartPasswordReset initiates the forgot password flow (sends email/saves code).
	// FIX: Updated return signature to include the generated code string for local debugging.
	StartPasswordReset(ctx context.Context, req domain.ForgotPasswordRequest) (string, error)

	// ResetPassword validates the code and updates the password.
	ResetPassword(ctx context.Context, req domain.ResetPasswordRequest) (*domain.AuthResponse, error)
}