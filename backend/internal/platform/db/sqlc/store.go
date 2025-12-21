package sqlc

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// Store wraps sqlc Queries and provides transaction support
type Store struct {
	*pgxpool.Pool
}

// NewStore creates a new store from a connection pool
func NewStore(pool *pgxpool.Pool) *Store {
	return &Store{Pool: pool}
}

// WithTx executes a function within a transaction, using the provided Queries interface
func (s *Store) WithTx(ctx context.Context, fn func(tx pgx.Tx) error) error {
	tx, err := s.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	if err := fn(tx); err != nil {
		return err
	}

	return tx.Commit(ctx)
}

// Querier is an interface that sqlc-generated Queries implement
// This allows us to work with queries in a generic way
type Querier interface {
	// Queries can be used directly with pool or transaction
}

