package command

import (
	"context"
	"fmt"

	"github.com/english-coach/backend/internal/domain/user/port"
	"go.uber.org/zap"
)

// UpdateUserProfileUseCase handles updating user profile
type UpdateUserProfileUseCase struct {
	profileRepo port.UserProfileRepository
	logger      *zap.Logger
}

// NewUpdateUserProfileUseCase creates a new update user profile use case
func NewUpdateUserProfileUseCase(
	profileRepo port.UserProfileRepository,
	logger *zap.Logger,
) *UpdateUserProfileUseCase {
	return &UpdateUserProfileUseCase{
		profileRepo: profileRepo,
		logger:      logger,
	}
}

// UpdateUserProfileInput represents the input for updating user profile
type UpdateUserProfileInput struct {
	DisplayName *string `json:"display_name,omitempty"`
	AvatarURL   *string `json:"avatar_url,omitempty"`
	BirthDay    *string `json:"birth_day,omitempty"` // Format: YYYY-MM-DD
	Bio         *string `json:"bio,omitempty"`
}

// UpdateUserProfileOutput represents the output for updating user profile
type UpdateUserProfileOutput struct {
	UserID      int64   `json:"user_id"`
	DisplayName *string `json:"display_name,omitempty"`
	AvatarURL   *string `json:"avatar_url,omitempty"`
	BirthDay    *string `json:"birth_day,omitempty"`
	Bio         *string `json:"bio,omitempty"`
}

// Execute updates user profile
func (uc *UpdateUserProfileUseCase) Execute(ctx context.Context, userID int64, input UpdateUserProfileInput) (*UpdateUserProfileOutput, error) {
	profile, err := uc.profileRepo.Update(ctx, userID, input.DisplayName, input.AvatarURL, input.BirthDay, input.Bio)
	if err != nil {
		uc.logger.Error("failed to update user profile", zap.Error(err), zap.Int64("user_id", userID))
		return nil, fmt.Errorf("failed to update user profile: %w", err)
	}

	var birthDayStr *string
	if profile.BirthDay != nil {
		formatted := profile.BirthDay.Format("2006-01-02")
		birthDayStr = &formatted
	}

	return &UpdateUserProfileOutput{
		UserID:      profile.UserID,
		DisplayName: profile.DisplayName,
		AvatarURL:   profile.AvatarURL,
		BirthDay:    birthDayStr,
		Bio:         profile.Bio,
	}, nil
}
