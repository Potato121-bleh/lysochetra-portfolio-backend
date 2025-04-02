package mockUtil

import (
	"context"
	"profile-portfolio/internal/db"

	// "github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/stretchr/testify/mock"
)

type MockTx struct {
	mock.Mock
}

func (t *MockTx) Query(ctx context.Context, sql string, args ...interface{}) (db.RowsScanner, error) {
	mockArgs := t.Called(ctx, sql, args)
	return mockArgs.Get(0).(db.RowsScanner), mockArgs.Error(1)
}

func (t *MockTx) QueryRow(ctx context.Context, sql string, args ...interface{}) db.RowScanner {
	mockArgs := t.Called(ctx, sql, args)
	return mockArgs.Get(0).(db.RowScanner)
}

func (t *MockTx) Exec(ctx context.Context, sql string, args ...interface{}) (pgconn.CommandTag, error) {
	mockArgs := t.Called(ctx, sql, args)
	return mockArgs.Get(0).(pgconn.CommandTag), mockArgs.Error(1)
}

func (t *MockTx) Begin(ctx context.Context) (db.DatabaseTx, error) {
	mockArgs := t.Called(ctx)
	return mockArgs.Get(0).(db.DatabaseTx), mockArgs.Error(1)
}

func (t *MockTx) Rollback(ctx context.Context) error {
	mockArgs := t.Called(ctx)
	return mockArgs.Error(0)
}

func (t *MockTx) Commit(ctx context.Context) error {
	mockArgs := t.Called(ctx)
	return mockArgs.Error(0)
}
