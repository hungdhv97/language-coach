package errors

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

// GetUniqueConstraintField extracts the field name from a unique violation error
// This helps identify which field caused the conflict
func GetUniqueConstraintField(err error) string {
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		// PostgreSQL unique violation error detail usually contains the constraint name
		// Format: "Key (field_name)=(value) already exists."
		// We can extract the field from the constraint name or detail
		if pgErr.Detail != "" {
			// Try to extract field name from detail
			// This is a simple implementation - can be enhanced based on actual error format
			return pgErr.ConstraintName
		}
	}
	return ""
}

// MapPgError is a generic mapper for PostgreSQL errors in repositories that don't have domain-specific mappings
// For most cases, it returns the error as-is. This function can be extended to handle common PostgreSQL errors
// if needed in the future.
func MapPgError(err error) error {
	if err == nil {
		return nil
	}
	// For now, return error as-is
	// This can be extended to handle common PostgreSQL errors if needed
	return err
}
