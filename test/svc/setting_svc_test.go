package svcTest

import (
	"fmt"
	"profile-portfolio/internal/domain/model"
	"profile-portfolio/internal/domain/service"
	mockUtil "profile-portfolio/test/mock"
	testUtilTool "profile-portfolio/test/util"
	"testing"
	"time"

	"github.com/jackc/pgx/v5/pgconn"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestSettingSelectSvc(t *testing.T) {
	timeNow := time.Now()
	mockDB := new(mockUtil.MockDB)
	mockTx := new(mockUtil.MockTx)
	expectedRs := []model.SettingStruct{
		{
			SettingId:     3,
			Darkmode:      1,
			Sound:         1,
			Colorpalettes: 0,
			Font:          0,
			Language:      0,
		},
		{
			SettingId:     2,
			Darkmode:      0,
			Sound:         1,
			Colorpalettes: 0,
			Font:          1,
			Language:      0,
		},
	}
	mockRowContext := [][]interface{}{
		{
			3,
			1,
			1,
			0,
			0,
			0,
		},
		{
			2,
			0,
			1,
			0,
			1,
			0,
		},
	}
	mockRow := mockUtil.NewMockRow(mockRowContext)

	mockScanArg := testUtilTool.CountMockAnything(model.SettingStruct{})

	mockRow.On("Scan", mockScanArg...).Return(nil)

	for testCaseIndex, ele := range SelectTestCase {
		t.Run(
			ele.name,
			func(t *testing.T) {
				prepMockDB := ele.setup(t, mockRow, mockTx, mockDB, mockScanArg)
				userSvc := service.NewSettingService(prepMockDB)
				prepTx := ele.tx(mockRowContext, timeNow, mockScanArg)
				testRs, testErr := userSvc.Select(prepTx, "mytb", "", "")
				if testCaseIndex == 2 {
					require.NotNil(t, testErr)
				} else {
					require.Nil(t, testErr)
					for i := 0; i < len(expectedRs); i++ {
						assert.Equal(t, expectedRs[i], testRs[i])
					}
				}

			},
		)
	}
}

func TestSettingInsertSvc(t *testing.T) {

	t.Run(
		"Test INSERT Service with: provided Tx | Expected No Error",
		func(t *testing.T) {
			mockTx := new(mockUtil.MockTx)
			mockDB := new(mockUtil.MockDB)
			mockTx.On("Commit", mock.Anything).Return(nil)
			mockTx.On("Exec", mock.Anything, "INSERT INTO mytb ( username , password ) VALUES ( $1 , $2 )", []interface{}{"username1", "password1"}).Return(pgconn.NewCommandTag("INSERT 1"), nil)
			mockDB.On("Begin", mock.Anything).Return(mockTx, nil)

			userSvc := service.NewSettingService(mockDB)
			testRs := userSvc.Insert(mockTx, "mytb", []string{"username", "password"}, []string{"username1", "password1"})
			require.Nil(t, testRs)
		},
	)

	t.Run(
		"Test INSERT Service with: provided Tx | Expected Error",
		func(t *testing.T) {
			mockTx := new(mockUtil.MockTx)
			mockDB := new(mockUtil.MockDB)
			mockTx.On("Commit", mock.Anything).Return(nil)
			mockTx.On("Rollback", mock.Anything).Return(nil)
			mockTx.On(
				"Exec",
				mock.Anything,
				"INSERT INTO mytb ( username , password ) VALUES ( $1 , $2 )",
				[]interface{}{"username1", "password1"},
			).Return(
				pgconn.NewCommandTag("INSERT 0"),
				fmt.Errorf("failed to execute INSERT stmt"),
			)

			mockDB.On("Begin", mock.Anything).Return(mockTx, nil)

			userSvc := service.NewSettingService(mockDB)
			testRs := userSvc.Insert(mockTx, "mytb", []string{"username", "password"}, []string{"username1", "password1"})
			require.NotNil(t, testRs)
		},
	)

}

func TestSettingUpdateSvc(t *testing.T) {
	t.Run(
		"Test UPDATE Service with: provided Tx | No WHERE Clause | Expected No Error",
		func(t *testing.T) {
			mockTx := new(mockUtil.MockTx)
			mockDB := new(mockUtil.MockDB)
			mockTx.On("Commit", mock.Anything).Return(nil)
			mockTx.On("Exec", mock.Anything, "UPDATE mytb SET username = $1 , password = $2", []interface{}{"username1", "password1"}).Return(pgconn.NewCommandTag("INSERT 10"), nil)
			mockDB.On("Begin", mock.Anything).Return(mockTx, nil)

			userSvc := service.NewSettingService(mockDB)
			testRs := userSvc.Update(mockTx, "mytb", []string{"username", "password"}, []string{"username1", "password1"}, "", "")
			require.Nil(t, testRs)
		},
	)

	t.Run(
		"Test UPDATE Service with: provided Tx | with WHERE Clause | Expected Error",
		func(t *testing.T) {
			mockTx := new(mockUtil.MockTx)
			mockDB := new(mockUtil.MockDB)
			mockTx.On("Commit", mock.Anything).Return(nil)
			mockTx.On("Rollback", mock.Anything).Return(nil)
			mockTx.On(
				"Exec",
				mock.Anything,
				"UPDATE mytb SET username = $1 , password = $2 WHERE setting_id = $3",
				[]interface{}{"username1", "password1", "1"},
			).Return(
				pgconn.NewCommandTag("INSERT 1"),
				fmt.Errorf("failed to execute INSERT stmt"),
			)

			mockDB.On("Begin", mock.Anything).Return(mockTx, nil)

			userSvc := service.NewSettingService(mockDB)
			testRs := userSvc.Update(mockTx, "mytb", []string{"username", "password"}, []string{"username1", "password1"}, "setting_id", "1")
			require.NotNil(t, testRs)
		},
	)
}

func TestSettingDeleteSvc(t *testing.T) {
	t.Run(
		"Test DELETE Service with: provided Tx | No WHERE Clause | Expected No Error",
		func(t *testing.T) {
			mockTx := new(mockUtil.MockTx)
			mockDB := new(mockUtil.MockDB)
			mockTx.On("Commit", mock.Anything).Return(nil)
			mockTx.On("Exec", mock.Anything, "DELETE FROM mytb WHERE user_id = $1", []interface{}{"1"}).Return(pgconn.NewCommandTag("DELETE 10"), nil)
			mockDB.On("Begin", mock.Anything).Return(mockTx, nil)

			userSvc := service.NewSettingService(mockDB)
			testRs := userSvc.Delete(mockTx, "mytb", "user_id", "1")
			require.Nil(t, testRs)
		},
	)

	t.Run(
		"Test DELETE Service with: provided Tx | with WHERE Clause | Expected Error",
		func(t *testing.T) {
			mockTx := new(mockUtil.MockTx)
			mockDB := new(mockUtil.MockDB)
			mockTx.On("Commit", mock.Anything).Return(nil)
			mockTx.On("Rollback", mock.Anything).Return(nil)
			mockTx.On(
				"Exec",
				mock.Anything,
				"DELETE FROM mytb WHERE user_id = $1",
				[]interface{}{"1"},
			).Return(
				pgconn.NewCommandTag("DELETE 0"),
				fmt.Errorf("failed to execute DELETE stmt"),
			)

			mockDB.On("Begin", mock.Anything).Return(mockTx, nil)

			userSvc := service.NewSettingService(mockDB)
			testRs := userSvc.Delete(mockTx, "mytb", "user_id", "1")
			require.NotNil(t, testRs)
		},
	)
}
