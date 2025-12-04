package pgxtype

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Tx interface {
	Commit(ctx context.Context) error
	Rollback(ctx context.Context) error
}

type TxWrapper struct {
	Tx pgx.Tx
}

func (t *TxWrapper) Commit(ctx context.Context) error {
	return t.Tx.Commit(ctx)
}

func (t *TxWrapper) Rollback(ctx context.Context) error {
	return t.Tx.Rollback(ctx)
}

func WithTx[T any](ctx context.Context, pool *pgxpool.Pool, wtx interface{ WithTx(pgx.Tx) *T }) (Tx, *T, error) {
	tx, err := pool.Begin(ctx)
	if err != nil {
		return nil, nil, fmt.Errorf("pgxpool.Begin: %w", err)
	}

	txQueries := wtx.WithTx(tx)

	return &TxWrapper{Tx: tx}, txQueries, nil
}
