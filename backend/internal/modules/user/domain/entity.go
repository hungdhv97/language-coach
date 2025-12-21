package domain

import "time"

// User represents a user in the system
type User struct {
	ID           int64     `json:"id"`
	Email        *string   `json:"email,omitempty"`
	Username     *string   `json:"username,omitempty"`
	PasswordHash string    `json:"-"` // Never expose password hash in JSON
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
	IsActive     bool      `json:"is_active"`
}

// UserProfile represents extended user profile information
type UserProfile struct {
	UserID      int64      `json:"user_id"`
	DisplayName *string    `json:"display_name,omitempty"`
	AvatarURL   *string    `json:"avatar_url,omitempty"`
	BirthDay    *time.Time `json:"birth_day,omitempty"`
	Bio         *string    `json:"bio,omitempty"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
}

