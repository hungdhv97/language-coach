package command

import (
	"context"
	"errors"
	"fmt"

	"github.com/english-coach/backend/internal/domain/user/model"
	"github.com/english-coach/backend/internal/domain/user/port"
	"github.com/english-coach/backend/internal/infrastructure/auth"
	"github.com/english-coach/backend/internal/infrastructure/crypto"
	"github.com/english-coach/backend/internal/infrastructure/repository/common"
	"go.uber.org/zap"
)

var (
	ErrInvalidCredentials = errors.New("invalid email/username or password")
	ErrUserInactive       = errors.New("user account is inactive")
)

// LoginUseCase handles user login
type LoginUseCase struct {
	userRepo   port.UserRepository
	jwtManager *auth.JWTManager
	logger     *zap.Logger
}

// NewLoginUseCase creates a new login use case
func NewLoginUseCase(
	userRepo port.UserRepository,
	jwtManager *auth.JWTManager,
	logger *zap.Logger,
) *LoginUseCase {
	return &LoginUseCase{
		userRepo:   userRepo,
		jwtManager: jwtManager,
		logger:     logger,
	}
}

// LoginInput represents the input for user login
type LoginInput struct {
	Email    *string `json:"email,omitempty"`
	Username *string `json:"username,omitempty"`
	Password string  `json:"password"`
}

// LoginOutput represents the output for user login
type LoginOutput struct {
	Token    string  `json:"token"`
	UserID   int64   `json:"user_id"`
	Email    *string `json:"email,omitempty"`
	Username *string `json:"username,omitempty"`
}

// Execute executes the user login
func (uc *LoginUseCase) Execute(ctx context.Context, input LoginInput) (*LoginOutput, error) {
	// Find user by email or username
	var user *model.User
	var err error

	if input.Email != nil && *input.Email != "" {
		user, err = uc.userRepo.FindByEmail(ctx, *input.Email)
	} else if input.Username != nil && *input.Username != "" {
		user, err = uc.userRepo.FindByUsername(ctx, *input.Username)
	} else {
		return nil, ErrInvalidCredentials
	}

	if err != nil {
		if common.IsNotFound(err) {
			return nil, ErrInvalidCredentials
		}
		uc.logger.Error("failed to find user", zap.Error(err))
		return nil, fmt.Errorf("failed to find user: %w", err)
	}

	if user == nil {
		return nil, ErrInvalidCredentials
	}

	// Check if user is active
	if !user.IsActive {
		return nil, ErrUserInactive
	}

	// Verify password
	if !crypto.CheckPasswordHash(input.Password, user.PasswordHash) {
		return nil, ErrInvalidCredentials
	}

	// Generate JWT token
	username := ""
	if user.Username != nil {
		username = *user.Username
	} else if user.Email != nil {
		username = *user.Email
	}

	token, err := uc.jwtManager.GenerateToken(user.ID, username)
	if err != nil {
		uc.logger.Error("failed to generate token", zap.Error(err))
		return nil, fmt.Errorf("failed to generate token: %w", err)
	}

	return &LoginOutput{
		Token:    token,
		UserID:   user.ID,
		Email:    user.Email,
		Username: user.Username,
	}, nil
}
