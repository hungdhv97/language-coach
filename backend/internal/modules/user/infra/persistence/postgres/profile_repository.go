package user

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5/pgtype"

	"github.com/english-coach/backend/internal/modules/user/domain"
	db "github.com/english-coach/backend/internal/platform/db/sqlc/gen/user"
	"github.com/english-coach/backend/internal/shared/errors"
)

// userProfileRepository implements domain.UserProfileRepository
type userProfileRepository struct {
	*UserRepository
}

// Create creates a new user profile
func (r *userProfileRepository) Create(ctx context.Context, userID int64, displayName *string, avatarURL *string, birthDay *string, bio *string) (*domain.UserProfile, error) {
	var displayNamePg pgtype.Text
	if displayName != nil && *displayName != "" {
		displayNamePg = pgtype.Text{String: *displayName, Valid: true}
	}

	var avatarURLPg pgtype.Text
	if avatarURL != nil && *avatarURL != "" {
		avatarURLPg = pgtype.Text{String: *avatarURL, Valid: true}
	}

	var birthDayPg pgtype.Date
	if birthDay != nil && *birthDay != "" {
		parsed, err := time.Parse("2006-01-02", *birthDay)
		if err == nil {
			birthDayPg = pgtype.Date{Time: parsed, Valid: true}
		}
	}

	var bioPg pgtype.Text
	if bio != nil && *bio != "" {
		bioPg = pgtype.Text{String: *bio, Valid: true}
	}

	row, err := r.queries.CreateUserProfile(ctx, db.CreateUserProfileParams{
		UserID:      userID,
		DisplayName: displayNamePg,
		AvatarUrl:   avatarURLPg,
		BirthDay:    birthDayPg,
		Bio:         bioPg,
	})
	if err != nil {
		return nil, errors.MapPgError(err)
	}

	return mapDBProfileToModel(&row), nil
}

// GetByUserID returns a user profile by user ID
func (r *userProfileRepository) GetByUserID(ctx context.Context, userID int64) (*domain.UserProfile, error) {
	row, err := r.queries.GetUserProfile(ctx, userID)
	if err != nil {
		return nil, errors.MapPgError(err)
	}

	return mapDBProfileToModel(&row), nil
}

// Update updates a user profile
func (r *userProfileRepository) Update(ctx context.Context, userID int64, displayName *string, avatarURL *string, birthDay *string, bio *string) (*domain.UserProfile, error) {
	var displayNamePg pgtype.Text
	if displayName != nil && *displayName != "" {
		displayNamePg = pgtype.Text{String: *displayName, Valid: true}
	}

	var avatarURLPg pgtype.Text
	if avatarURL != nil && *avatarURL != "" {
		avatarURLPg = pgtype.Text{String: *avatarURL, Valid: true}
	}

	var birthDayPg pgtype.Date
	if birthDay != nil && *birthDay != "" {
		parsed, err := time.Parse("2006-01-02", *birthDay)
		if err == nil {
			birthDayPg = pgtype.Date{Time: parsed, Valid: true}
		}
	}

	var bioPg pgtype.Text
	if bio != nil && *bio != "" {
		bioPg = pgtype.Text{String: *bio, Valid: true}
	}

	row, err := r.queries.UpdateUserProfile(ctx, db.UpdateUserProfileParams{
		UserID:      userID,
		DisplayName: displayNamePg,
		AvatarUrl:   avatarURLPg,
		BirthDay:    birthDayPg,
		Bio:         bioPg,
	})
	if err != nil {
		return nil, errors.MapPgError(err)
	}

	return mapDBProfileToModel(&row), nil
}

// mapDBProfileToModel maps sqlc generated UserProfile to domain model
func mapDBProfileToModel(row *db.UserProfile) *domain.UserProfile {
	var displayName *string
	var avatarURL *string
	var birthDay *time.Time
	var bio *string

	if row.DisplayName.Valid {
		displayName = &row.DisplayName.String
	}
	if row.AvatarUrl.Valid {
		avatarURL = &row.AvatarUrl.String
	}
	if row.BirthDay.Valid {
		birthDay = &row.BirthDay.Time
	}
	if row.Bio.Valid {
		bio = &row.Bio.String
	}

	return &domain.UserProfile{
		UserID:      row.UserID,
		DisplayName: displayName,
		AvatarURL:   avatarURL,
		BirthDay:    birthDay,
		Bio:         bio,
		CreatedAt:   row.CreatedAt.Time,
		UpdatedAt:   row.UpdatedAt.Time,
	}
}
