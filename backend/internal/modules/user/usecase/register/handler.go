package register

import (
	"context"

	"github.com/english-coach/backend/internal/modules/user/domain"
	"github.com/english-coach/backend/internal/shared/auth"
	"github.com/english-coach/backend/internal/shared/errors"
	"github.com/english-coach/backend/internal/shared/logger"
)

// Handler handles user registration
type Handler struct {
	userRepo domain.UserRepository
	logger   logger.ILogger
}

// NewHandler creates a new register user handler
func NewHandler(
	userRepo domain.UserRepository,
	logger logger.ILogger,
) *Handler {
	return &Handler{
		userRepo: userRepo,
		logger:   logger,
	}
}

// Execute executes the user registration
func (h *Handler) Execute(ctx context.Context, input Input) (*Output, error) {
	// Validate input
	if (input.Email == nil || *input.Email == "") && (input.Username == nil || *input.Username == "") {
		return nil, domain.ErrEmailRequired
	}

	if len(input.Password) < 6 {
		return nil, domain.ErrInvalidPassword
	}

	// Check if email already exists
	if input.Email != nil && *input.Email != "" {
		exists, err := h.userRepo.CheckEmailExists(ctx, *input.Email)
		if err != nil {
			h.logger.Error("failed to check if email exists",
				logger.Error(err),
				logger.String("email", *input.Email),
			)
			return nil, errors.WrapError(err, "failed to check if email exists")
		}
		if exists {
			return nil, domain.ErrEmailExists
		}
	}

	// Check if username already exists
	if input.Username != nil && *input.Username != "" {
		exists, err := h.userRepo.CheckUsernameExists(ctx, *input.Username)
		if err != nil {
			h.logger.Error("failed to check if username exists",
				logger.Error(err),
				logger.String("username", *input.Username),
			)
			return nil, errors.WrapError(err, "failed to check if username exists")
		}
		if exists {
			return nil, domain.ErrUsernameExists
		}
	}

	// Hash password
	passwordHash, err := auth.HashPassword(input.Password)
	if err != nil {
		h.logger.Error("failed to hash password", logger.Error(err))
		return nil, errors.WrapError(err, "failed to hash password")
	}

	// Create user
	user, err := h.userRepo.Create(ctx, input.Email, input.Username, passwordHash)
	if err != nil {
		h.logger.Error("failed to create user", logger.Error(err))
		return nil, errors.WrapError(err, "failed to create user")
	}

	return &Output{
		UserID:   user.ID,
		Email:    user.Email,
		Username: user.Username,
	}, nil
}
