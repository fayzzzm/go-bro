package postgres

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// queryOne is a generic helper to execute a query and collect a single row into a struct.
func queryOne[T any](ctx context.Context, pool *pgxpool.Pool, query string, args ...any) (*T, error) {
	rows, err := pool.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	// CollectOneRow handles closing the rows
	val, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[T])

	if err != nil {
		return nil, err
	}

	return &val, nil
}

// queryRows is a generic helper to execute a query and collect multiple rows into a slice of structs.
func queryRows[T any](ctx context.Context, pool *pgxpool.Pool, query string, args ...any) ([]T, error) {
	rows, err := pool.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}

	// CollectRows handles closing the rows
	return pgx.CollectRows(rows, pgx.RowToStructByName[T])
}
