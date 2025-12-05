package port

import (
	"context"

	"github.com/english-coach/backend/internal/domain/user/model"
)

// UserRepository defines operations for user data access
type UserRepository interface {
	// Create creates a new user
	Create(ctx context.Context, email *string, username *string, passwordHash string) (*model.User, error)
	// FindByID returns a user by ID
	FindByID(ctx context.Context, id int64) (*model.User, error)
	// FindByEmail returns a user by email
	FindByEmail(ctx context.Context, email string) (*model.User, error)
	// FindByUsername returns a user by username
	FindByUsername(ctx context.Context, username string) (*model.User, error)
	// UpdatePassword updates a user's password
	UpdatePassword(ctx context.Context, id int64, passwordHash string) error
	// UpdateActiveStatus updates a user's active status
	UpdateActiveStatus(ctx context.Context, id int64, isActive bool) error
	// CheckEmailExists checks if an email already exists
	CheckEmailExists(ctx context.Context, email string) (bool, error)
	// CheckUsernameExists checks if a username already exists
	CheckUsernameExists(ctx context.Context, username string) (bool, error)
}

// UserProfileRepository defines operations for user profile data access
type UserProfileRepository interface {
	// Create creates a new user profile
	Create(ctx context.Context, userID int64, displayName *string, avatarURL *string, birthDay *string, bio *string) (*model.UserProfile, error)
	// GetByUserID returns a user profile by user ID
	GetByUserID(ctx context.Context, userID int64) (*model.UserProfile, error)
	// Update updates a user profile
	Update(ctx context.Context, userID int64, displayName *string, avatarURL *string, birthDay *string, bio *string) (*model.UserProfile, error)
}
