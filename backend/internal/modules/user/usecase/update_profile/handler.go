package update_profile

import (
	"context"

	"github.com/english-coach/backend/internal/modules/user/domain"
	"github.com/english-coach/backend/internal/shared/errors"
	"github.com/english-coach/backend/internal/shared/logger"
)

// Handler handles updating user profile
type Handler struct {
	profileRepo domain.UserProfileRepository
	logger      logger.ILogger
}

// NewHandler creates a new update user profile handler
func NewHandler(
	profileRepo domain.UserProfileRepository,
	logger logger.ILogger,
) *Handler {
	return &Handler{
		profileRepo: profileRepo,
		logger:      logger,
	}
}

// Execute updates user profile
func (h *Handler) Execute(ctx context.Context, userID int64, input Input) (*Output, error) {
	profile, err := h.profileRepo.Update(ctx, userID, input.DisplayName, input.AvatarURL, input.BirthDay, input.Bio)
	if err != nil {
		h.logger.Error("failed to update user profile", logger.Error(err), logger.Int64("user_id", userID))
		return nil, errors.WrapError(err, "failed to update user profile")
	}

	var birthDayStr *string
	if profile.BirthDay != nil {
		formatted := profile.BirthDay.Format("2006-01-02")
		birthDayStr = &formatted
	}

	return &Output{
		UserID:      profile.UserID,
		DisplayName: profile.DisplayName,
		AvatarURL:   profile.AvatarURL,
		BirthDay:    birthDayStr,
		Bio:         profile.Bio,
	}, nil
}
