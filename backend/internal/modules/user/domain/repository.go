package domain

import (
	"context"
)

// UserRepository defines operations for user data access
type UserRepository interface {
	// Create creates a new user
	Create(ctx context.Context, email *string, username *string, passwordHash string) (*User, error)
	// FindUserByID returns a user by ID
	FindUserByID(ctx context.Context, id int64) (*User, error)
	// FindUserByEmail returns a user by email
	FindUserByEmail(ctx context.Context, email string) (*User, error)
	// FindUserByUsername returns a user by username
	FindUserByUsername(ctx context.Context, username string) (*User, error)
	// UpdatePassword updates a user's password
	UpdatePassword(ctx context.Context, id int64, passwordHash string) error
	// UpdateActiveStatus updates a user's active status
	UpdateActiveStatus(ctx context.Context, id int64, isActive bool) error
	// ExistsEmail checks if an email already exists
	ExistsEmail(ctx context.Context, email string) (bool, error)
	// ExistsUsername checks if a username already exists
	ExistsUsername(ctx context.Context, username string) (bool, error)
}

// UserProfileRepository defines operations for user profile data access
type UserProfileRepository interface {
	// Create creates a new user profile
	Create(ctx context.Context, userID int64, displayName *string, avatarURL *string, birthDay *string, bio *string) (*UserProfile, error)
	// FindUserProfileByUserID returns a user profile by user ID
	FindUserProfileByUserID(ctx context.Context, userID int64) (*UserProfile, error)
	// Update updates a user profile
	Update(ctx context.Context, userID int64, displayName *string, avatarURL *string, birthDay *string, bio *string) (*UserProfile, error)
}
