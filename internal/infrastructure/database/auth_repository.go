package database

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/mitcheltastic/ManproBackend/internal/core/domain"
	"github.com/mitcheltastic/ManproBackend/internal/core/ports"
)

// AuthRepository implements the ports.AuthRepository interface for Postgres (Supabase).
type AuthRepository struct {
	DB *sql.DB
}

// NewAuthRepository creates a new instance of the AuthRepository.
func NewAuthRepository(db *sql.DB) ports.AuthRepository {
	return &AuthRepository{DB: db}
}

// CreateUser saves a new user record to the 'users' table.
func (r *AuthRepository) CreateUser(ctx context.Context, user domain.User) error {
	query := `
		INSERT INTO users (id, name, email, hashed_password, is_verified, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
	`
	_, err := r.DB.ExecContext(
		ctx,
		query,
		user.ID,
		user.Name,
		user.Email,
		user.HashedPassword,
		user.IsVerified,
		user.CreatedAt,
		user.UpdatedAt,
	)
	if err != nil {
		// In a real app, you would check for unique constraint violations here (e.g., duplicate email)
		return err 
	}
	return nil
}

// GetUserByEmail retrieves a user by their email address.
func (r *AuthRepository) GetUserByEmail(ctx context.Context, email string) (*domain.User, error) {
	query := `
		SELECT id, name, email, hashed_password, is_verified, created_at, updated_at 
		FROM users WHERE email = $1
	`
	user := &domain.User{}
	err := r.DB.QueryRowContext(ctx, query, email).Scan(
		&user.ID,
		&user.Name,
		&user.Email,
		&user.HashedPassword,
		&user.IsVerified,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil // User not found
		}
		return nil, err
	}
	return user, nil
}

// UpdateUserPassword updates the user's password hash in the database.
func (r *AuthRepository) UpdateUserPassword(ctx context.Context, userID string, newHashedPassword string) error {
	query := `
		UPDATE users SET hashed_password = $1, updated_at = $2 WHERE id = $3
	`
	_, err := r.DB.ExecContext(ctx, query, newHashedPassword, time.Now(), userID)
	return err
}

// --- Password Reset Logic (Requires a temporary 'password_reset_tokens' table) ---

// CreatePasswordResetCode saves a unique 6-digit code linked to a user/email.
func (r *AuthRepository) CreatePasswordResetCode(ctx context.Context, email string, code string) error {
	// IMPORTANT: You will need to create a table in Supabase like:
	// CREATE TABLE password_reset_tokens (email TEXT PRIMARY KEY, code TEXT, expires_at TIMESTAMP WITH TIME ZONE);
	
	// Upsert (insert or update) the code, setting it to expire in 15 minutes.
	query := `
		INSERT INTO password_reset_tokens (email, code, expires_at)
		VALUES ($1, $2, $3)
		ON CONFLICT (email) DO UPDATE 
		SET code = EXCLUDED.code, expires_at = EXCLUDED.expires_at
	`
	expiresAt := time.Now().Add(15 * time.Minute)
	_, err := r.DB.ExecContext(ctx, query, email, code, expiresAt)
	return err
}

// VerifyPasswordResetCode checks if the code is valid and not expired.
func (r *AuthRepository) VerifyPasswordResetCode(ctx context.Context, email string, code string) error {
	query := `
		SELECT code, expires_at FROM password_reset_tokens WHERE email = $1
	`
	var storedCode string
	var expiresAt time.Time
	
	err := r.DB.QueryRowContext(ctx, query, email).Scan(&storedCode, &expiresAt)
	
	if errors.Is(err, sql.ErrNoRows) {
		return errors.New("reset code not found for this email")
	}
	if err != nil {
		return err
	}
	
	if storedCode != code {
		return errors.New("invalid verification code")
	}

	if expiresAt.Before(time.Now()) {
		return errors.New("verification code expired")
	}

	return nil // Code is valid and not expired
}

// DeletePasswordResetCode removes the code after a successful reset.
func (r *AuthRepository) DeletePasswordResetCode(ctx context.Context, email string) error {
	query := `DELETE FROM password_reset_tokens WHERE email = $1`
	_, err := r.DB.ExecContext(ctx, query, email)
	return err
}