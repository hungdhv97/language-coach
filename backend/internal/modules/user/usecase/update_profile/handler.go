package update_profile

import (
	"context"

	"github.com/english-coach/backend/internal/modules/user/domain"
	sharederrors "github.com/english-coach/backend/internal/shared/errors"
)

// Handler handles updating user profile
type Handler struct {
	profileRepo domain.UserProfileRepository
}

// NewHandler creates a new update user profile handler
func NewHandler(
	profileRepo domain.UserProfileRepository,
) *Handler {
	return &Handler{
		profileRepo: profileRepo,
	}
}

// Execute updates user profile
func (h *Handler) Execute(ctx context.Context, userID int64, input UpdateProfileInput) (*UpdateProfileOutput, error) {
	profile, err := h.profileRepo.Update(ctx, userID, input.DisplayName, input.AvatarURL, input.BirthDay, input.Bio)
	if err != nil {
		// Map domain error to AppError
		return nil, sharederrors.MapDomainErrorToAppError(err)
	}

	var birthDayStr *string
	if profile.BirthDay != nil {
		formatted := profile.BirthDay.Format("2006-01-02")
		birthDayStr = &formatted
	}

	return &UpdateProfileOutput{
		UserID:      profile.UserID,
		DisplayName: profile.DisplayName,
		AvatarURL:   profile.AvatarURL,
		BirthDay:    birthDayStr,
		Bio:         profile.Bio,
	}, nil
}
