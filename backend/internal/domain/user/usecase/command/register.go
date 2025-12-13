package command

import (
	"context"
	"errors"
	"fmt"

	"github.com/english-coach/backend/internal/domain/user/port"
	"github.com/english-coach/backend/internal/infrastructure/crypto"
	"go.uber.org/zap"
)

var (
	ErrEmailRequired   = errors.New("email or username is required")
	ErrEmailExists     = errors.New("email already exists")
	ErrUsernameExists  = errors.New("username already exists")
	ErrInvalidPassword = errors.New("password must be at least 6 characters")
)

// RegisterUserUseCase handles user registration
type RegisterUserUseCase struct {
	userRepo port.UserRepository
	logger   *zap.Logger
}

// NewRegisterUserUseCase creates a new register user use case
func NewRegisterUserUseCase(
	userRepo port.UserRepository,
	logger *zap.Logger,
) *RegisterUserUseCase {
	return &RegisterUserUseCase{
		userRepo: userRepo,
		logger:   logger,
	}
}

// RegisterUserInput represents the input for user registration
type RegisterUserInput struct {
	DisplayName *string `json:"display_name,omitempty"`
	Email       *string `json:"email,omitempty"`
	Username    *string `json:"username,omitempty"`
	Password    string  `json:"password"`
}

// RegisterUserOutput represents the output for user registration
type RegisterUserOutput struct {
	UserID   int64   `json:"user_id"`
	Email    *string `json:"email,omitempty"`
	Username *string `json:"username,omitempty"`
}

// Execute executes the user registration
func (uc *RegisterUserUseCase) Execute(ctx context.Context, input RegisterUserInput) (*RegisterUserOutput, error) {
	// Validate input
	if (input.Email == nil || *input.Email == "") && (input.Username == nil || *input.Username == "") {
		return nil, ErrEmailRequired
	}

	if len(input.Password) < 6 {
		return nil, ErrInvalidPassword
	}

	// Check if email already exists
	if input.Email != nil && *input.Email != "" {
		exists, err := uc.userRepo.CheckEmailExists(ctx, *input.Email)
		if err != nil {
			uc.logger.Error("failed to check if email exists",
				zap.Error(err),
				zap.String("email", *input.Email),
			)
			return nil, fmt.Errorf("failed to check if email exists: %w", err)
		}
		if exists {
			return nil, ErrEmailExists
		}
	}

	// Check if username already exists
	if input.Username != nil && *input.Username != "" {
		exists, err := uc.userRepo.CheckUsernameExists(ctx, *input.Username)
		if err != nil {
			uc.logger.Error("failed to check if username exists",
				zap.Error(err),
				zap.Stringp("username", input.Username),
			)
			return nil, fmt.Errorf("failed to check if username exists: %w", err)
		}
		if exists {
			return nil, ErrUsernameExists
		}
	}

	// Hash password
	passwordHash, err := crypto.HashPassword(input.Password)
	if err != nil {
		uc.logger.Error("failed to hash password", zap.Error(err))
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	// Create user
	user, err := uc.userRepo.Create(ctx, input.Email, input.Username, passwordHash)
	if err != nil {
		uc.logger.Error("failed to create user", zap.Error(err))
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	return &RegisterUserOutput{
		UserID:   user.ID,
		Email:    user.Email,
		Username: user.Username,
	}, nil
}
