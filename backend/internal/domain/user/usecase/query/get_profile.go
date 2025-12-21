package query

import (
	"context"

	"github.com/english-coach/backend/internal/modules/user/domain"
	"github.com/english-coach/backend/internal/shared/errors"
	"go.uber.org/zap"
)

// GetUserProfileUseCase handles getting user profile
type GetUserProfileUseCase struct {
	profileRepo domain.UserProfileRepository
	logger      *zap.Logger
}

// NewGetUserProfileUseCase creates a new get user profile use case
func NewGetUserProfileUseCase(
	profileRepo domain.UserProfileRepository,
	logger *zap.Logger,
) *GetUserProfileUseCase {
	return &GetUserProfileUseCase{
		profileRepo: profileRepo,
		logger:      logger,
	}
}

// Execute gets user profile
func (uc *GetUserProfileUseCase) Execute(ctx context.Context, userID int64) (*domain.UserProfile, error) {
	profile, err := uc.profileRepo.GetByUserID(ctx, userID)
	if err != nil {
		uc.logger.Error("failed to get user profile", zap.Error(err), zap.Int64("user_id", userID))
		return nil, errors.WrapError(err, "failed to get user profile")
	}

	if profile == nil {
		return nil, domain.ErrProfileNotFound
	}

	return profile, nil
}
