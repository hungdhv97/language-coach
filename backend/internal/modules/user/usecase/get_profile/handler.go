package get_profile

import (
	"context"

	"github.com/english-coach/backend/internal/modules/user/domain"
	sharederrors "github.com/english-coach/backend/internal/shared/errors"
)

// Handler handles getting user profile
type Handler struct {
	profileRepo domain.UserProfileRepository
}

// NewHandler creates a new get user profile handler
func NewHandler(
	profileRepo domain.UserProfileRepository,
) *Handler {
	return &Handler{
		profileRepo: profileRepo,
	}
}

// Execute gets user profile
func (h *Handler) Execute(ctx context.Context, input GetProfileInput) (*GetProfileOutput, error) {
	profile, err := h.profileRepo.GetByUserID(ctx, input.UserID)
	if err != nil {
		// Map domain error to AppError
		return nil, sharederrors.MapDomainErrorToAppError(err)
	}

	if profile == nil {
		return nil, sharederrors.MapDomainErrorToAppError(domain.ErrProfileNotFound)
	}

	var birthDayStr *string
	if profile.BirthDay != nil {
		formatted := profile.BirthDay.Format("2006-01-02")
		birthDayStr = &formatted
	}

	return &GetProfileOutput{
		UserID:      profile.UserID,
		DisplayName: profile.DisplayName,
		AvatarURL:   profile.AvatarURL,
		BirthDay:    birthDayStr,
		Bio:         profile.Bio,
	}, nil
}
