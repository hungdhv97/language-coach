package command

import (
	"context"

	"github.com/english-coach/backend/internal/modules/user/domain"
	"github.com/english-coach/backend/internal/shared/auth"
	"github.com/english-coach/backend/internal/shared/errors"
	"go.uber.org/zap"
)

// LoginUseCase handles user login
type LoginUseCase struct {
	userRepo   domain.UserRepository
	jwtManager *auth.JWTManager
	logger     *zap.Logger
}

// NewLoginUseCase creates a new login use case
func NewLoginUseCase(
	userRepo domain.UserRepository,
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
	var user *domain.User
	var err error

	if input.Email != nil && *input.Email != "" {
		user, err = uc.userRepo.FindByEmail(ctx, *input.Email)
	} else if input.Username != nil && *input.Username != "" {
		user, err = uc.userRepo.FindByUsername(ctx, *input.Username)
	} else {
		return nil, domain.ErrInvalidCredentials
	}

	if err != nil {
		if errors.IsNotFound(err) {
			return nil, domain.ErrInvalidCredentials
		}
		uc.logger.Error("failed to find user", zap.Error(err))
		return nil, errors.WrapError(err, "failed to find user")
	}

	if user == nil {
		return nil, domain.ErrInvalidCredentials
	}

	// Check if user is active
	if !user.IsActive {
		return nil, domain.ErrUserInactive
	}

	// Verify password
	if !auth.CheckPasswordHash(input.Password, user.PasswordHash) {
		return nil, domain.ErrInvalidCredentials
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
		return nil, errors.WrapError(err, "failed to generate token")
	}

	return &LoginOutput{
		Token:    token,
		UserID:   user.ID,
		Email:    user.Email,
		Username: user.Username,
	}, nil
}
