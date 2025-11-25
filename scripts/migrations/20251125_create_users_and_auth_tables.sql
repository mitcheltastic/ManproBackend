-- +goose Up
-- Creates the core user and password reset tables for the authentication system.

-- Table 1: users
CREATE TABLE users (
    -- Primary key, using a standard UUID generator for uniqueness and security
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    
    name VARCHAR(255) NOT NULL,
    email VARCHAR(255) UNIQUE NOT NULL,

    -- Stores the bcrypt hash of the password
    hashed_password TEXT NOT NULL,
    
    -- Used for email verification flow (if implemented later)
    is_verified BOOLEAN DEFAULT FALSE,
    
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW() NOT NULL,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW() NOT NULL
);

-- Table 2: password_reset_tokens
CREATE TABLE password_reset_tokens (
    -- Email is the foreign key and the primary key for simple lookup
    email TEXT PRIMARY KEY REFERENCES users(email) ON DELETE CASCADE,
    
    -- The 6-digit code sent to the user
    code CHAR(6) NOT NULL,
    
    -- Time when the code becomes invalid
    expires_at TIMESTAMP WITH TIME ZONE NOT NULL
);

-- +goose Down
-- Reverts the schema changes (used for rolling back a migration)
DROP TABLE password_reset_tokens;
DROP TABLE users;