package login

import (
	"context"
	"errors"

	"github.com/english-coach/backend/internal/modules/user/domain"
	"github.com/english-coach/backend/internal/shared/auth"
	sharederrors "github.com/english-coach/backend/internal/shared/errors"
)

// Handler handles user login
type Handler struct {
	userRepo   domain.UserRepository
	jwtManager *auth.JWTManager
}

// NewHandler creates a new login handler
func NewHandler(
	userRepo domain.UserRepository,
	jwtManager *auth.JWTManager,
) *Handler {
	return &Handler{
		userRepo:   userRepo,
		jwtManager: jwtManager,
	}
}

// Execute executes the user login
func (h *Handler) Execute(ctx context.Context, input LoginInput) (*LoginOutput, error) {
	// Find user by email or username
	var user *domain.User
	var err error

	if input.Email != nil && *input.Email != "" {
		user, err = h.userRepo.FindByEmail(ctx, *input.Email)
	} else if input.Username != nil && *input.Username != "" {
		user, err = h.userRepo.FindByUsername(ctx, *input.Username)
	} else {
		return nil, sharederrors.MapDomainErrorToAppError(domain.ErrInvalidCredentials)
	}

	if err != nil {
		// User not found -> invalid credentials (security: don't reveal if user exists)
		if errors.Is(err, domain.ErrUserNotFound) {
			return nil, sharederrors.MapDomainErrorToAppError(domain.ErrInvalidCredentials)
		}
		// Map domain error to AppError
		return nil, sharederrors.MapDomainErrorToAppError(err)
	}

	if user == nil {
		return nil, sharederrors.MapDomainErrorToAppError(domain.ErrInvalidCredentials)
	}

	// Check if user is active
	if !user.IsActive {
		return nil, sharederrors.MapDomainErrorToAppError(domain.ErrUserInactive)
	}

	// Verify password
	if !auth.CheckPasswordHash(input.Password, user.PasswordHash) {
		return nil, sharederrors.MapDomainErrorToAppError(domain.ErrInvalidCredentials)
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
		// Unexpected error from auth layer
		return nil, sharederrors.WrapUnexpectedError(err, "failed to generate token")
	}

	return &LoginOutput{
		Token:    token,
		UserID:   user.ID,
		Email:    user.Email,
		Username: user.Username,
	}, nil
}
