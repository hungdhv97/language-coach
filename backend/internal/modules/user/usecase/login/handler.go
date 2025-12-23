package login

import (
	"context"

	"github.com/english-coach/backend/internal/modules/user/domain"
	"github.com/english-coach/backend/internal/shared/auth"
	"github.com/english-coach/backend/internal/shared/errors"
	"github.com/english-coach/backend/internal/shared/logger"
)

// Handler handles user login
type Handler struct {
	userRepo   domain.UserRepository
	jwtManager *auth.JWTManager
	logger     logger.ILogger
}

// NewHandler creates a new login handler
func NewHandler(
	userRepo domain.UserRepository,
	jwtManager *auth.JWTManager,
	logger logger.ILogger,
) *Handler {
	return &Handler{
		userRepo:   userRepo,
		jwtManager: jwtManager,
		logger:     logger,
	}
}

// Execute executes the user login
func (h *Handler) Execute(ctx context.Context, input Input) (*Output, error) {
	// Find user by email or username
	var user *domain.User
	var err error

	if input.Email != nil && *input.Email != "" {
		user, err = h.userRepo.FindByEmail(ctx, *input.Email)
	} else if input.Username != nil && *input.Username != "" {
		user, err = h.userRepo.FindByUsername(ctx, *input.Username)
	} else {
		return nil, domain.ErrInvalidCredentials
	}

	if err != nil {
		// Check if it's a not found error (can be checked via errors package if needed)
		return nil, domain.ErrInvalidCredentials
		h.logger.Error("failed to find user", logger.Error(err))
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

	token, err := h.jwtManager.GenerateToken(user.ID, username)
	if err != nil {
		h.logger.Error("failed to generate token", logger.Error(err))
		return nil, errors.WrapError(err, "failed to generate token")
	}

	return &Output{
		Token:    token,
		UserID:   user.ID,
		Email:    user.Email,
		Username: user.Username,
	}, nil
}
