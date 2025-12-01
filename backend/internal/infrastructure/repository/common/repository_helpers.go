package common

import (
	"errors"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

// IsUniqueViolation checks if the error is a unique constraint violation
func IsUniqueViolation(err error) bool {
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		return pgErr.Code == "23505" // unique_violation
	}
	return false
}

// IsNotFound checks if the error is a "not found" error (pgx.ErrNoRows)
func IsNotFound(err error) bool {
	return errors.Is(err, pgx.ErrNoRows)
}

// MapPgError maps PostgreSQL errors to domain errors
// This can be extended to map specific error codes to domain-specific errors
func MapPgError(err error) error {
	if IsNotFound(err) {
		return err // Return as-is, let domain layer handle
	}
	if IsUniqueViolation(err) {
		return err // Return as-is, let domain layer handle
	}
	return err
}

