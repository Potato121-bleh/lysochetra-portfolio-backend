package db

import (
	"context"

	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Database interface {
	Query(ctx context.Context, sql string, args ...interface{}) (RowsScanner, error)
	QueryRow(ctx context.Context, sql string, args ...interface{}) RowScanner
	Exec(ctx context.Context, sql string, args ...interface{}) (pgconn.CommandTag, error)
	Begin(ctx context.Context) (DatabaseTx, error)
	Ping(ctx context.Context) error
	Close()
}

func NewPgxDBAdapter(db *pgxpool.Pool) *PgxDBAdapter {
	return &PgxDBAdapter{
		db: db,
	}
}

// create databaseTx adapter
type PgxDBAdapter struct {
	db *pgxpool.Pool
}

func (d *PgxDBAdapter) Query(ctx context.Context, sql string, args ...interface{}) (RowsScanner, error) {
	return d.db.Query(ctx, sql, args...)
}

func (d *PgxDBAdapter) QueryRow(ctx context.Context, sql string, args ...interface{}) RowScanner {
	return d.db.QueryRow(ctx, sql, args...)
}

func (d *PgxDBAdapter) Exec(ctx context.Context, sql string, args ...interface{}) (pgconn.CommandTag, error) {
	return d.db.Exec(ctx, sql, args...)
}

func (d *PgxDBAdapter) Begin(ctx context.Context) (DatabaseTx, error) {
	newTx, txErr := d.db.Begin(ctx)
	if txErr != nil {
		return nil, txErr
	}
	return &PgxTxAdapter{
		tx: newTx,
	}, nil
}

func (d *PgxDBAdapter) Ping(ctx context.Context) error {
	return d.db.Ping(ctx)
}

func (d *PgxDBAdapter) Close() {
	return
}

// func (d *PgxDBAdapter) BeginTx(ctx context.Context, txOptions pgx.TxOptions) (DatabaseTx, error) {
// 	newTx, txErr := d.db.BeginTx(ctx, )
// 	if txErr != nil {
// 		return nil, txErr
// 	}
// 	return &PgxTxAdapter{
// 		tx: newTx,
// 	}, nil
// }

// func (d *PgxDBAdapter) Rollback(ctx context.Context) error {
// 	return d.db.Rollback(ctx)
// }

// func (d *PgxDBAdapter) Commit(ctx context.Context) error {
// 	return d.tx.Commit(ctx)
// }
