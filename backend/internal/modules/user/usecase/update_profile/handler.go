package update_profile

import (
	"context"

	"github.com/english-coach/backend/internal/modules/user/domain"
	"github.com/english-coach/backend/internal/shared/errors"
	"go.uber.org/zap"
)

// Handler handles updating user profile
type Handler struct {
	profileRepo domain.UserProfileRepository
	logger      *zap.Logger
}

// NewHandler creates a new update user profile handler
func NewHandler(
	profileRepo domain.UserProfileRepository,
	logger *zap.Logger,
) *Handler {
	return &Handler{
		profileRepo: profileRepo,
		logger:      logger,
	}
}

// Input represents the input for updating user profile
type Input struct {
	DisplayName *string `json:"display_name,omitempty"`
	AvatarURL   *string `json:"avatar_url,omitempty"`
	BirthDay    *string `json:"birth_day,omitempty"` // Format: YYYY-MM-DD
	Bio         *string `json:"bio,omitempty"`
}

// Output represents the output for updating user profile
type Output struct {
	UserID      int64   `json:"user_id"`
	DisplayName *string `json:"display_name,omitempty"`
	AvatarURL   *string `json:"avatar_url,omitempty"`
	BirthDay    *string `json:"birth_day,omitempty"`
	Bio         *string `json:"bio,omitempty"`
}

// Execute updates user profile
func (h *Handler) Execute(ctx context.Context, userID int64, input Input) (*Output, error) {
	profile, err := h.profileRepo.Update(ctx, userID, input.DisplayName, input.AvatarURL, input.BirthDay, input.Bio)
	if err != nil {
		h.logger.Error("failed to update user profile", zap.Error(err), zap.Int64("user_id", userID))
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
