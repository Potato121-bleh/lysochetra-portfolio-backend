package db

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

type DatabaseTx interface {
	Query(ctx context.Context, sql string, args ...any) (RowsScanner, error)
	QueryRow(ctx context.Context, sql string, args ...any) RowScanner
	Exec(ctx context.Context, sql string, args ...any) (pgconn.CommandTag, error)
	Begin(ctx context.Context) (DatabaseTx, error)
	Rollback(ctx context.Context) error
	Commit(ctx context.Context) error
}

func NewPgxTxAdapter(tx pgx.Tx) *PgxTxAdapter {
	return &PgxTxAdapter{
		tx: tx,
	}
}

// create databaseTx adapter
type PgxTxAdapter struct {
	tx pgx.Tx
}

func (d *PgxTxAdapter) Query(ctx context.Context, sql string, args ...any) (RowsScanner, error) {
	return d.tx.Query(ctx, sql, args...)
}

func (d *PgxTxAdapter) QueryRow(ctx context.Context, sql string, args ...any) RowScanner {
	return d.tx.QueryRow(ctx, sql, args...)
}

func (d *PgxTxAdapter) Exec(ctx context.Context, sql string, args ...any) (pgconn.CommandTag, error) {
	return d.tx.Exec(ctx, sql, args...)
}

func (d *PgxTxAdapter) Begin(ctx context.Context) (DatabaseTx, error) {
	newTx, txErr := d.tx.Begin(ctx)
	if txErr != nil {
		return nil, txErr
	}
	return &PgxTxAdapter{
		tx: newTx,
	}, nil
}

func (d *PgxTxAdapter) Rollback(ctx context.Context) error {
	return d.tx.Rollback(ctx)
}

func (d *PgxTxAdapter) Commit(ctx context.Context) error {
	return d.tx.Commit(ctx)
}
