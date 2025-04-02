package repoTest

import (
	"profile-portfolio/internal/domain/model"
	"profile-portfolio/internal/domain/repository"
	mockUtil "profile-portfolio/test/mock"
	testUtilTool "profile-portfolio/test/util"
	"testing"
	"time"

	"github.com/jackc/pgx/v5/pgconn"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestUserSelectRepo(t *testing.T) {
	timeNow := time.Now()
	userRepo := repository.NewUserRepository()
	mockTx := new(mockUtil.MockTx)
	mockRow := mockUtil.NewMockRow(
		[][]interface{}{
			{
				1,
				"username1",
				"nickname1",
				"password1",
				timeNow,
				1,
			},
		})

	expectRs := model.UserData{
		Id:             1,
		Username:       "username1",
		Nickname:       "nickname1",
		Password:       "password1",
		RegisteredDate: timeNow,
		SettingId:      1,
	}

	mockScanArg := testUtilTool.CountMockAnything(model.UserData{})
	mockRow.On("Scan", mockScanArg...).Return(nil)
	mockTx.On("Query", mock.Anything, "SELECT * FROM mytb WHERE user_id = $1", []interface{}{"1"}).Return(mockRow, nil)
	testRs, testStatus := userRepo.SqlSelect(mockTx, "mytb", "user_id", "1")
	t.Run(
		"Compare expected Value with Result",
		func(t *testing.T) {
			require.Nil(t, testStatus)
			require.Equal(t, expectRs, testRs[0])
		},
	)

}

func TestUserInsertRepo(t *testing.T) {
	userRepo := repository.NewUserRepository()

	mockTx := new(mockUtil.MockTx)
	mockTx.On("Exec", mock.Anything, "INSERT INTO mytb ( username , password ) VALUES ( $1 , $2 )", []interface{}{"username1", "password1"}).Return(pgconn.NewCommandTag("INSERT 1"), nil)
	testRs := userRepo.SqlInsert(mockTx, "mytb", []string{"username", "password"}, []string{"username1", "password1"})
	t.Run(
		"Comparison between expected Row affect of Insert Repository",
		func(t *testing.T) {
			assert.Nil(t, testRs)
		},
	)
}

func TestUserUpdateRepo(t *testing.T) {
	userRepo := repository.NewUserRepository()

	t.Run(
		"Test UPDATE repository without WHERE Clause",
		func(t *testing.T) {
			mockTx := new(mockUtil.MockTx)
			mockTx.On("Exec", mock.Anything, "UPDATE mytb SET username = $1 , password = $2", []interface{}{"username1", "password1"}).Return(pgconn.NewCommandTag("UPDATE 10"), nil)

			testRs := userRepo.SqlUpdate(mockTx, "mytb", []string{"username", "password"}, []string{"username1", "password1"}, "", "")
			assert.Nil(t, testRs)
		},
	)

	t.Run(
		"Test UPDATE repository With WHERE Clause",
		func(t *testing.T) {
			mockTx := new(mockUtil.MockTx)
			mockTx.On(
				"Exec",
				mock.Anything,
				"UPDATE mytb SET username_kilo = $1 , password_kilo = $2 WHERE setting_id = $3",
				[]interface{}{"username_kilo1", "password_kilo1", "2"},
			).Return(pgconn.NewCommandTag("UPDATE 1"), nil)
			testRs := userRepo.SqlUpdate(mockTx, "mytb", []string{"username_kilo", "password_kilo"}, []string{"username_kilo1", "password_kilo1"}, "setting_id", "2")
			assert.Nil(t, testRs)
		},
	)

}

func TestUserDeleteRepo(t *testing.T) {
	userRepo := repository.NewUserRepository()

	t.Run(
		"Test DELETE repository without WHERE Clause",
		func(t *testing.T) {
			// Expected stmt: DELETE FROM mytb
			mockTx := new(mockUtil.MockTx)
			mockTx.On("Exec", mock.Anything, "DELETE FROM mytb", "").Return(pgconn.NewCommandTag("DELETE 100"), nil)
			testRs := userRepo.SqlDelete(mockTx, "mytb", "", "")
			assert.NotNil(t, testRs)
		},
	)

	t.Run(
		"Test DELETE repository with WHERE Clause",
		func(t *testing.T) {
			// Expected stmt: DELETE FROM mytb
			mockTx := new(mockUtil.MockTx)
			mockTx.On("Exec", mock.Anything, "DELETE FROM mytb WHERE username = $1", []interface{}{"1"}).Return(pgconn.NewCommandTag("DELETE 1"), nil)
			testRs := userRepo.SqlDelete(mockTx, "mytb", "username", "1")
			assert.Nil(t, testRs)
		},
	)

}
