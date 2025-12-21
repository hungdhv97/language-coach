package db

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// WithTx executes a function within a database transaction.
// If the function returns an error, the transaction is rolled back.
// Otherwise, it is committed.
func WithTx(ctx context.Context, pool *pgxpool.Pool, fn func(tx pgx.Tx) error) error {
	tx, err := pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	if err := fn(tx); err != nil {
		return err
	}

	return tx.Commit(ctx)
}

// BeginTx starts a new transaction and returns it along with a commit/rollback function.
func BeginTx(ctx context.Context, pool *pgxpool.Pool) (pgx.Tx, func() error, func() error, error) {
	tx, err := pool.Begin(ctx)
	if err != nil {
		return nil, nil, nil, err
	}

	commit := func() error {
		return tx.Commit(ctx)
	}

	rollback := func() error {
		return tx.Rollback(ctx)
	}

	return tx, commit, rollback, nil
}

