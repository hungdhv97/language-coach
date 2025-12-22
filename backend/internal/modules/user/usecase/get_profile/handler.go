package get_profile

import (
	"context"

	"github.com/english-coach/backend/internal/modules/user/domain"
	"github.com/english-coach/backend/internal/shared/errors"
	"go.uber.org/zap"
)

// Handler handles getting user profile
type Handler struct {
	profileRepo domain.UserProfileRepository
	logger      *zap.Logger
}

// NewHandler creates a new get user profile handler
func NewHandler(
	profileRepo domain.UserProfileRepository,
	logger *zap.Logger,
) *Handler {
	return &Handler{
		profileRepo: profileRepo,
		logger:      logger,
	}
}

// Execute gets user profile
func (h *Handler) Execute(ctx context.Context, input Input) (*Output, error) {
	profile, err := h.profileRepo.GetByUserID(ctx, input.UserID)
	if err != nil {
		h.logger.Error("failed to get user profile", zap.Error(err), zap.Int64("user_id", input.UserID))
		return nil, errors.WrapError(err, "failed to get user profile")
	}

	if profile == nil {
		return nil, domain.ErrProfileNotFound
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
