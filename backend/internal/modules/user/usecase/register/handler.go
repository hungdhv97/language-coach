package register

import (
	"context"

	"github.com/english-coach/backend/internal/modules/user/domain"
	"github.com/english-coach/backend/internal/shared/auth"
	sharederrors "github.com/english-coach/backend/internal/shared/errors"
)

// Handler handles user registration
type Handler struct {
	userRepo domain.UserRepository
}

// NewHandler creates a new register user handler
func NewHandler(
	userRepo domain.UserRepository,
) *Handler {
	return &Handler{
		userRepo: userRepo,
	}
}

// Execute executes the user registration
func (h *Handler) Execute(ctx context.Context, input RegisterInput) (*RegisterOutput, error) {
	// Validate input
	if (input.Email == nil || *input.Email == "") && (input.Username == nil || *input.Username == "") {
		return nil, sharederrors.MapDomainErrorToAppError(domain.ErrEmailRequired)
	}

	if len(input.Password) < 6 {
		return nil, sharederrors.MapDomainErrorToAppError(domain.ErrInvalidPassword)
	}

	// Check if email already exists
	if input.Email != nil && *input.Email != "" {
		exists, err := h.userRepo.ExistsEmail(ctx, *input.Email)
		if err != nil {
			// Map domain error to AppError
			return nil, sharederrors.MapDomainErrorToAppError(err)
		}
		if exists {
			return nil, sharederrors.MapDomainErrorToAppError(domain.ErrEmailExists)
		}
	}

	// Check if username already exists
	if input.Username != nil && *input.Username != "" {
		exists, err := h.userRepo.ExistsUsername(ctx, *input.Username)
		if err != nil {
			// Map domain error to AppError
			return nil, sharederrors.MapDomainErrorToAppError(err)
		}
		if exists {
			return nil, sharederrors.MapDomainErrorToAppError(domain.ErrUsernameExists)
		}
	}

	// Hash password
	passwordHash, err := auth.HashPassword(input.Password)
	if err != nil {
		return nil, sharederrors.MapDomainErrorToAppError(err)
	}

	// Create user
	user, err := h.userRepo.Create(ctx, input.Email, input.Username, passwordHash)
	if err != nil {
		// Map domain error to AppError (MapDomainErrorToAppError handles all cases)
		return nil, sharederrors.MapDomainErrorToAppError(err)
	}

	return &RegisterOutput{
		UserID:   user.ID,
		Email:    user.Email,
		Username: user.Username,
	}, nil
}
