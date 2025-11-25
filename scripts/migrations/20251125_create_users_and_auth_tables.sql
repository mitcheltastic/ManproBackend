-- +goose Up
-- Creates the core user and password reset tables for the authentication system.

-- Table 1: users
CREATE TABLE users (
    -- Core Auth/Security Fields
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(255) NOT NULL,
    email VARCHAR(255) UNIQUE NOT NULL,
    hashed_password TEXT NOT NULL,
    is_verified BOOLEAN DEFAULT FALSE,

    -- Intertwine Profile Fields (All nullable until profile is set)
    nickname VARCHAR(50) NULL,
    gender VARCHAR(10) NULL, -- 'M', 'F', 'Other'
    birthdate DATE NULL, -- DATE is preferred for birthdate over TIMESTAMP
    college VARCHAR(100) NULL,
    faculty VARCHAR(100) NULL,
    major VARCHAR(100) NULL,
    year INT NULL,
    mbti CHAR(4) NULL, -- e.g., 'INFP'
    blood_type CHAR(2) NULL, -- e.g., 'AB'
    profile_picture_url TEXT NULL, -- URL to the primary image
    
    -- Array types for multiple values
    gallery_picture_urls TEXT[] NULL, -- Array of image URLs (max 5)
    hobby TEXT[] NULL, -- Array of strings for hobbies (max 3)

    -- Metadata Fields
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW() NOT NULL,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW() NOT NULL
);

-- Table 2: password_reset_tokens
CREATE TABLE password_reset_tokens (
    email TEXT PRIMARY KEY REFERENCES users(email) ON DELETE CASCADE,
    code CHAR(6) NOT NULL,
    expires_at TIMESTAMP WITH TIME ZONE NOT NULL
);

-- +goose Down
-- FIX: Use CASCADE on the parent table (users) for robustness.
DROP TABLE users CASCADE;
DROP TABLE password_reset_tokens CASCADE;