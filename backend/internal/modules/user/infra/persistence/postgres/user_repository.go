package user

import (
	"context"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/english-coach/backend/internal/modules/user/domain"
	db "github.com/english-coach/backend/internal/platform/db/sqlc/gen/user"
	sharederrors "github.com/english-coach/backend/internal/shared/errors"
)

// UserRepository implements user repository interfaces using sqlc
type UserRepository struct {
	pool    *pgxpool.Pool
	queries *db.Queries
}

// NewUserRepository creates a new user repository
func NewUserRepository(pool *pgxpool.Pool) *UserRepository {
	return &UserRepository{
		pool:    pool,
		queries: db.New(pool),
	}
}

// UserRepository returns a UserRepository implementation
func (r *UserRepository) UserRepository() domain.UserRepository {
	return &userRepository{
		UserRepository: r,
	}
}

// UserProfileRepository returns a UserProfileRepository implementation
func (r *UserRepository) UserProfileRepository() domain.UserProfileRepository {
	return &userProfileRepository{
		UserRepository: r,
	}
}

// userRepository implements domain.UserRepository
type userRepository struct {
	*UserRepository
}

// Create creates a new user
func (r *userRepository) Create(ctx context.Context, email *string, username *string, passwordHash string) (*domain.User, error) {
	var emailPg pgtype.Text
	if email != nil && *email != "" {
		emailPg = pgtype.Text{String: *email, Valid: true}
	}

	var usernamePg pgtype.Text
	if username != nil && *username != "" {
		usernamePg = pgtype.Text{String: *username, Valid: true}
	}

	row, err := r.queries.CreateUser(ctx, db.CreateUserParams{
		Email:        emailPg,
		Username:     usernamePg,
		PasswordHash: pgtype.Text{String: passwordHash, Valid: true},
		IsActive:     pgtype.Bool{Bool: true, Valid: true},
	})
	if err != nil {
		return nil, sharederrors.MapUserRepositoryError(err, "Create")
	}

	return mapDBUserToModel(&row), nil
}

// FindUserByID returns a user by ID
func (r *userRepository) FindUserByID(ctx context.Context, id int64) (*domain.User, error) {
	row, err := r.queries.FindUserByID(ctx, id)
	if err != nil {
		return nil, sharederrors.MapUserRepositoryError(err, "FindUserByID")
	}

	return mapDBUserToModel(&row), nil
}

// FindUserByEmail returns a user by email
func (r *userRepository) FindUserByEmail(ctx context.Context, email string) (*domain.User, error) {
	row, err := r.queries.FindUserByEmail(ctx, pgtype.Text{String: email, Valid: true})
	if err != nil {
		return nil, sharederrors.MapUserRepositoryError(err, "FindUserByEmail")
	}

	return mapDBUserToModel(&row), nil
}

// FindUserByUsername returns a user by username
func (r *userRepository) FindUserByUsername(ctx context.Context, username string) (*domain.User, error) {
	row, err := r.queries.FindUserByUsername(ctx, pgtype.Text{String: username, Valid: true})
	if err != nil {
		return nil, sharederrors.MapUserRepositoryError(err, "FindUserByUsername")
	}

	return mapDBUserToModel(&row), nil
}

// UpdatePassword updates a user's password
func (r *userRepository) UpdatePassword(ctx context.Context, id int64, passwordHash string) error {
	err := r.queries.UpdateUserPassword(ctx, db.UpdateUserPasswordParams{
		ID:           id,
		PasswordHash: pgtype.Text{String: passwordHash, Valid: true},
	})
	return sharederrors.MapUserRepositoryError(err, "UpdatePassword")
}

// UpdateActiveStatus updates a user's active status
func (r *userRepository) UpdateActiveStatus(ctx context.Context, id int64, isActive bool) error {
	err := r.queries.UpdateUserActiveStatus(ctx, db.UpdateUserActiveStatusParams{
		ID:       id,
		IsActive: pgtype.Bool{Bool: isActive, Valid: true},
	})
	return sharederrors.MapUserRepositoryError(err, "UpdateActiveStatus")
}

// ExistsEmail checks if an email already exists
func (r *userRepository) ExistsEmail(ctx context.Context, email string) (bool, error) {
	result, err := r.queries.CheckEmailExists(ctx, pgtype.Text{String: email, Valid: true})
	if err != nil {
		return false, sharederrors.MapUserRepositoryError(err, "ExistsEmail")
	}
	return result, nil
}

// ExistsUsername checks if a username already exists
func (r *userRepository) ExistsUsername(ctx context.Context, username string) (bool, error) {
	result, err := r.queries.CheckUsernameExists(ctx, pgtype.Text{String: username, Valid: true})
	if err != nil {
		return false, sharederrors.MapUserRepositoryError(err, "ExistsUsername")
	}
	return result, nil
}

// mapDBUserToModel maps sqlc generated User to domain model
func mapDBUserToModel(row *db.User) *domain.User {
	var email *string
	var username *string

	if row.Email.Valid {
		email = &row.Email.String
	}
	if row.Username.Valid {
		username = &row.Username.String
	}

	return &domain.User{
		ID:           row.ID,
		Email:        email,
		Username:     username,
		PasswordHash: row.PasswordHash.String,
		CreatedAt:    row.CreatedAt.Time,
		UpdatedAt:    row.UpdatedAt.Time,
		IsActive:     row.IsActive.Bool,
	}
}
