package service

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/google/uuid"
	"github.com/mitcheltastic/ManproBackend/internal/core/domain"
	"github.com/mitcheltastic/ManproBackend/internal/core/ports"
	"github.com/mitcheltastic/ManproBackend/internal/pkg/email" // CRITICAL IMPORT: Ensure this is present
	"github.com/mitcheltastic/ManproBackend/internal/pkg/security"
)

// AuthService is the concrete implementation of the ports.AuthService interface.
type AuthService struct {
	AuthRepo ports.AuthRepository
	JWTService security.JWTService 
	EmailSender email.Sender // CRITICAL: Email Sender dependency must be here
}

// NewAuthService creates a new instance of the AuthService.
// CRITICAL: The third argument (emailSender) must be accepted here.
func NewAuthService(authRepo ports.AuthRepository, jwtService security.JWTService, emailSender email.Sender) ports.AuthService {
	return &AuthService{
		AuthRepo: authRepo,
		JWTService: jwtService,
		EmailSender: emailSender,
	}
}

// Register handles user registration, including password hashing and storage.
func (s *AuthService) Register(ctx context.Context, req domain.RegisterRequest) (*domain.AuthResponse, error) {
	// 1. Check if user already exists
	existingUser, err := s.AuthRepo.GetUserByEmail(ctx, req.Email)
	if err != nil {
		return nil, fmt.Errorf("repository error during lookup: %w", err)
	}
	if existingUser != nil {
		return nil, errors.New("a user with this email already exists")
	}

	// 2. Hash the password
	hashedPassword, err := security.HashPassword(req.Password)
	if err != nil {
		return nil, errors.New("failed to hash password")
	}

	// 3. Create the new User domain model
	newUser := domain.User{
		ID:        uuid.New(),
		Name:      req.Name,
		Email:     req.Email,
		HashedPassword: hashedPassword,
		IsVerified:    false, // Typically requires email confirmation
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	// 4. Save the user to the database
	if err := s.AuthRepo.CreateUser(ctx, newUser); err != nil {
		return nil, fmt.Errorf("failed to save new user: %w", err)
	}

	// 5. Generate JWT token
	token, err := s.JWTService.GenerateToken(newUser.ID, newUser.Email)
	if err != nil {
		return nil, fmt.Errorf("failed to generate auth token: %w", err)
	}

	// 6. Return success response
	return &domain.AuthResponse{
		UserID: newUser.ID,
		Email:  newUser.Email,
		Name:   newUser.Name,
		Token:  token,
	}, nil
}

// Login verifies user credentials and issues an authentication token.
func (s *AuthService) Login(ctx context.Context, req domain.LoginRequest) (*domain.AuthResponse, error) {
	// 1. Retrieve user by email
	user, err := s.AuthRepo.GetUserByEmail(ctx, req.Email)
	if err != nil {
		return nil, errors.New("internal error during login")
	}
	if user == nil {
		return nil, errors.New("invalid credentials") // Use generic message for security
	}

	// 2. Compare the stored hash with the provided password
	if err := security.CheckPasswordHash(req.Password, user.HashedPassword); err != nil {
		return nil, errors.New("invalid credentials")
	}
	
	// 3. Generate JWT token
	token, err := s.JWTService.GenerateToken(user.ID, user.Email)
	if err != nil {
		return nil, fmt.Errorf("failed to generate auth token: %w", err)
	}

	// 4. Return success response
	return &domain.AuthResponse{
		UserID: user.ID,
		Email:  user.Email,
		Name:   user.Name,
		Token:  token,
	}, nil
}

// StartPasswordReset initiates the forgot password flow by generating and saving a reset code.
func (s *AuthService) StartPasswordReset(ctx context.Context, req domain.ForgotPasswordRequest) (string, error) {
	// 1. Check if user exists 
	user, err := s.AuthRepo.GetUserByEmail(ctx, req.Email)
	if err != nil {
		return "", errors.New("internal server error")
	}
	if user == nil {
		// IMPORTANT: For security, we return success even if the user doesn't exist.
		log.Printf("Password reset requested for non-existent email: %s", req.Email)
		return "", nil 
	}
	
	// 2. Generate a secure, 6-digit random code
	code := security.GenerateNumericCode(6)
	
	// 3. Save the code and its expiration time to the repository (password_reset_tokens table)
	if err := s.AuthRepo.CreatePasswordResetCode(ctx, req.Email, code); err != nil {
		return "", fmt.Errorf("failed to save reset code: %w", err)
	}

	// 4. CRITICAL FIX: Call the actual email sender here
	if err := s.EmailSender.SendPasswordResetCode(req.Email, code); err != nil {
		// Log the error but return success to avoid leaking internal email failures
		log.Printf("ERROR: Failed to send reset code email to %s: %v", req.Email, err)
		return "", nil // Return empty code string and nil error to satisfy the handler's requirement for security
	}
	
	return code, nil // Return the generated code for testing
}

// ResetPassword validates the code and updates the user's password.
func (s *AuthService) ResetPassword(ctx context.Context, req domain.ResetPasswordRequest) (*domain.AuthResponse, error) {
	// 1. Check if the reset code is valid and not expired
	if err := s.AuthRepo.VerifyPasswordResetCode(ctx, req.Email, req.Code); err != nil {
		return nil, err // Returns specific errors like "code expired" or "invalid code"
	}
	
	// 2. Retrieve user to get the User ID
	user, err := s.AuthRepo.GetUserByEmail(ctx, req.Email)
	if err != nil {
		return nil, errors.New("internal error retrieving user")
	}
	if user == nil {
		return nil, errors.New("user not found after code verification")
	}

	// 3. Hash the new password
	newHashedPassword, err := security.HashPassword(req.NewPassword)
	if err != nil {
		return nil, errors.New("failed to hash new password")
	}

	// 4. Update the user's password in the database
	if err := s.AuthRepo.UpdateUserPassword(ctx, user.ID.String(), newHashedPassword); err != nil {
		return nil, fmt.Errorf("failed to update password: %w", err)
	}

	// 5. Clean up: Delete the reset code
	if err := s.AuthRepo.DeletePasswordResetCode(ctx, req.Email); err != nil {
		log.Printf("Warning: Failed to delete reset code for %s: %v", req.Email, err)
	}

	// 6. Generate a new JWT token for the user
	token, err := s.JWTService.GenerateToken(user.ID, user.Email)
	if err != nil {
		return nil, fmt.Errorf("failed to generate auth token after reset: %w", err)
	}


	// 7. Return success response
	return &domain.AuthResponse{
		UserID: user.ID,
		Email:  user.Email,
		Name:   user.Name,
		Token:  token,
	}, nil
}