package mockUtil

import (
	"context"
	"profile-portfolio/internal/db"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/stretchr/testify/mock"
)

type MockDB struct {
	mock.Mock
}

func (d *MockDB) Query(ctx context.Context, sql string, args ...interface{}) (db.RowsScanner, error) {
	mockArgs := d.Called(ctx, sql, args)
	return mockArgs.Get(0).(db.RowsScanner), mockArgs.Error(1)
}

func (d *MockDB) QueryRow(ctx context.Context, sql string, args ...interface{}) db.RowScanner {
	mockArgs := d.Called(ctx, sql, args)
	return mockArgs.Get(0).(db.RowScanner)
}

func (d *MockDB) Exec(ctx context.Context, sql string, args ...interface{}) (pgconn.CommandTag, error) {
	mockArgs := d.Called(ctx, sql, args)
	return mockArgs.Get(0).(pgconn.CommandTag), mockArgs.Error(1)
}

func (d *MockDB) Begin(ctx context.Context) (db.DatabaseTx, error) {
	mockArgs := d.Called(ctx)
	return mockArgs.Get(0).(db.DatabaseTx), mockArgs.Error(1)
}

func (d *MockDB) BeginTx(ctx context.Context, txOptions pgx.TxOptions) (db.DatabaseTx, error) {
	mockArgs := d.Called(ctx, txOptions)
	return mockArgs.Get(0).(db.DatabaseTx), mockArgs.Error(1)
}

func (d *MockDB) Ping(ctx context.Context) error {
	mockArgs := d.Called(ctx)
	return mockArgs.Error(1)
}

func (d *MockDB) Close() {
	d.Called()
}
