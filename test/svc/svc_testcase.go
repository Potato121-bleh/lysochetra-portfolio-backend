package svcTest

import (
	"fmt"
	"profile-portfolio/internal/db"
	mockUtil "profile-portfolio/test/mock"
	"testing"
	"time"

	"github.com/stretchr/testify/mock"
)

type SelectTest struct {
	name  string
	setup func(t *testing.T, mockRow *mockUtil.MockRow, mockTx *mockUtil.MockTx, mockDB *mockUtil.MockDB, mockScanArg []interface{}) *mockUtil.MockDB
	tx    func(rowContext [][]interface{}, timeNow time.Time, mockScanArg []interface{}) db.DatabaseTx
}

var SelectTestCase = []SelectTest{
	{
		name: "Test SELECT Service with: No provided Tx, Sucessfully commit",
		setup: func(t *testing.T, mockRow *mockUtil.MockRow, mockTx *mockUtil.MockTx, mockDB *mockUtil.MockDB, mockScanArg []interface{}) *mockUtil.MockDB {
			mockTx.On("Commit", mock.Anything).Return(nil)
			mockTx.On("Query", mock.Anything, "SELECT * FROM mytb", []interface{}{""}).Return(mockRow, nil)
			mockDB.On("Begin", mock.Anything).Return(mockTx, nil)

			return mockDB
		},
		tx: func(rowContext [][]interface{}, timeNow time.Time, mockScanArg []interface{}) db.DatabaseTx {
			return nil
		},
	},
	{
		name: "Test SELECT Service with: provided Tx, Sucessfully commit",
		setup: func(t *testing.T, mockRow *mockUtil.MockRow, mockTx *mockUtil.MockTx, mockDB *mockUtil.MockDB, mockScanArg []interface{}) *mockUtil.MockDB {
			mockDB.On("Begin", mock.Anything).Return(mockTx, nil)

			return mockDB
		},
		tx: func(rowContext [][]interface{}, timeNow time.Time, mockScanArg []interface{}) db.DatabaseTx {
			mockTx := new(mockUtil.MockTx)
			mockRow := mockUtil.NewMockRow(rowContext)
			mockRow.On("Scan", mockScanArg...).Return(nil)
			mockTx.On("Commit", mock.Anything).Return(nil)
			mockTx.On("Rollback", mock.Anything).Return(nil)
			mockTx.On("Query", mock.Anything, "SELECT * FROM mytb", []interface{}{""}).Return(mockRow, nil)

			return mockTx
		},
	},
	{
		name: "Test SELECT Service with: provided Tx, query failed",
		setup: func(t *testing.T, mockRow *mockUtil.MockRow, mockTx *mockUtil.MockTx, mockDB *mockUtil.MockDB, mockScanArg []interface{}) *mockUtil.MockDB {
			mockDB.On("Begin", mock.Anything).Return(mockTx, nil)

			return mockDB
		},
		tx: func(rowContext [][]interface{}, timeNow time.Time, mockScanArg []interface{}) db.DatabaseTx {
			mockTx := new(mockUtil.MockTx)
			mockRow := mockUtil.NewMockRow(rowContext)
			mockRow.On("Scan", mockScanArg...).Return(nil)
			mockTx.On("Commit", mock.Anything).Return(nil)
			mockTx.On("Rollback", mock.Anything).Return(nil)
			mockTx.On("Query", mock.Anything, "SELECT * FROM mytb", []interface{}{""}).Return(mockRow, fmt.Errorf("failed to perform query transaction."))

			return mockTx
		},
	},
}
