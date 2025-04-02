package repoTest

import (
	"fmt"
	"profile-portfolio/internal/domain/model"
	"profile-portfolio/internal/domain/repository"
	mockUtil "profile-portfolio/test/mock"
	testUtilTool "profile-portfolio/test/util"
	"testing"

	"github.com/jackc/pgx/v5/pgconn"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestSettingSelectRepo(t *testing.T) {
	settingRepo := repository.NewSettingRepository()
	mockTx := new(mockUtil.MockTx)
	mockRow := mockUtil.NewMockRow(
		[][]interface{}{
			{
				1,
				1,
				0,
				0,
				0,
				1,
			},
			{
				2,
				0,
				1,
				1,
				0,
				1,
			},
			{
				3,
				1,
				1,
				1,
				0,
				1,
			},
		})

	expectRs := []model.SettingStruct{
		{
			SettingId:     1,
			Darkmode:      1,
			Sound:         0,
			Colorpalettes: 0,
			Font:          0,
			Language:      1,
		},
		{
			SettingId:     2,
			Darkmode:      0,
			Sound:         1,
			Colorpalettes: 1,
			Font:          0,
			Language:      1,
		},
		{
			SettingId:     3,
			Darkmode:      1,
			Sound:         1,
			Colorpalettes: 1,
			Font:          0,
			Language:      1,
		},
	}

	mockScanArg := testUtilTool.CountMockAnything(model.UserData{})
	mockRow.On("Scan", mockScanArg...).Return(nil)
	mockTx.On("Query", mock.Anything, "SELECT * FROM mytb WHERE setting_id = $1", []interface{}{"1"}).Return(mockRow, nil)
	testRs, testStatus := settingRepo.SqlSelect(mockTx, "mytb", "setting_id", "1")
	t.Run(
		"Compare expected Value with Result in Setting Repo",
		func(t *testing.T) {
			require.Nil(t, testStatus)
			for i := 0; i < len(testRs); i++ {
				require.Equal(t, expectRs[i], testRs[i])
			}
		},
	)

}

func TestSettingInsertRepo(t *testing.T) {
	userRepo := repository.NewSettingRepository()

	t.Run(
		"Test INSERT Setting Repository WITH 2 Row Affected",
		func(t *testing.T) {
			mockTx := new(mockUtil.MockTx)
			mockTx.On("Exec", mock.Anything, "INSERT INTO mytb ( setting_id , language ) VALUES ( $1 , $2 )", []interface{}{"3", "1"}).Return(pgconn.NewCommandTag("INSERT 2"), nil)
			testRs := userRepo.SqlInsert(mockTx, "mytb", []string{"setting_id", "language"}, []string{"3", "1"})
			require.NotNil(t, testRs)
		},
	)

	t.Run(
		"Test INSERT Setting Repository WITH 1 Row Affected",
		func(t *testing.T) {
			mockTx := new(mockUtil.MockTx)
			mockTx.On("Exec", mock.Anything, "INSERT INTO mytb ( setting_id , language ) VALUES ( $1 , $2 )", []interface{}{"3", "1"}).Return(pgconn.NewCommandTag("INSERT 1"), fmt.Errorf("failed to execute INSERT stmt"))
			testRs := userRepo.SqlInsert(mockTx, "mytb", []string{"setting_id", "language"}, []string{"3", "1"})
			require.NotNil(t, testRs)
		},
	)
}

func TestSettingUpdateRepo(t *testing.T) {
	settingRepo := repository.NewSettingRepository()

	t.Run(
		"Test UPDATE repository WITHOUT WHERE Clause, Expected to return error",
		func(t *testing.T) {
			mockTx := new(mockUtil.MockTx)
			mockTx.On(
				"Exec",
				mock.Anything,
				"UPDATE mytb SET setting_id = $1 , language = $2",
				[]interface{}{"5", "0"},
			).Return(pgconn.NewCommandTag("UPDATE 10"), fmt.Errorf("failed to execute update stmt"))
			testRs := settingRepo.SqlUpdate(mockTx, "mytb", []string{"setting_id", "language"}, []string{"5", "0"}, "", "")
			require.NotNil(t, testRs)
		},
	)

	t.Run(
		"Test UPDATE repository WITH WHERE Clause",
		func(t *testing.T) {
			mockTx := new(mockUtil.MockTx)
			mockTx.On(
				"Exec",
				mock.Anything,
				"UPDATE mytb SET setting_id = $1 , language_kilo = $2 WHERE user_id = $3",
				[]interface{}{"5", "1", "1"},
			).Return(pgconn.NewCommandTag("UPDATE 10"), nil)
			testRs := settingRepo.SqlUpdate(mockTx, "mytb", []string{"setting_id", "language_kilo"}, []string{"5", "1"}, "user_id", "1")
			require.Nil(t, testRs)
		},
	)
}

func TestSettingDeleteRepo(t *testing.T) {
	settingRepo := repository.NewSettingRepository()

	t.Run(
		"Test DELETE repository without WHERE Clause, Expected Error from database (Exec)",
		func(t *testing.T) {
			mockTx := new(mockUtil.MockTx)
			mockTx.On("Exec", mock.Anything, "DELETE FROM mytb", "").Return(pgconn.NewCommandTag("DELETE 100"), fmt.Errorf("failed to execute DELETE stmt"))
			testRs := settingRepo.SqlDelete(mockTx, "mytb", "", "")
			require.NotNil(t, testRs)
		},
	)

	t.Run(
		"Test DELETE repository with WHERE Clause",
		func(t *testing.T) {
			mockTx := new(mockUtil.MockTx)
			mockTx.On("Exec", mock.Anything, "DELETE FROM mytb WHERE username = $1", []interface{}{"1"}).Return(pgconn.NewCommandTag("DELETE 1"), nil)
			testRs := settingRepo.SqlDelete(mockTx, "mytb", "username", "1")
			require.Nil(t, testRs)
		},
	)

}
