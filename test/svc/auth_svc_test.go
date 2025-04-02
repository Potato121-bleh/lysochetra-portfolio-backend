package svcTest

import (
	"fmt"
	"net/http"
	"profile-portfolio/internal/domain/model"
	"profile-portfolio/internal/domain/service"
	mockUtil "profile-portfolio/test/mock"
	"testing"
	"time"

	"github.com/jackc/pgx/v5/pgconn"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

type signUpTestCaseType struct {
	name                  string
	scan                  error
	rowContext            []interface{}
	ExecExpectRowAffected string
	IsExpectRsNil         bool
}

var signUpTestCase = []signUpTestCaseType{
	{
		name: "Test SignUp with: all require, Excepted No Error",
		scan: nil,
		rowContext: []interface{}{
			2,
		},
		ExecExpectRowAffected: "INSERT 1",
		IsExpectRsNil:         true,
	},
	{
		name:                  "Test SignUp with: No Row found on Scan, Excepted Error",
		scan:                  fmt.Errorf("no row found"),
		rowContext:            []interface{}{},
		ExecExpectRowAffected: "INSERT 1",
		IsExpectRsNil:         false,
	},
	{
		name: "Test SignUp with: AffectRow 0 Exec for INSERT userauth, Excepted Error",
		scan: nil,
		rowContext: []interface{}{
			2,
		},
		ExecExpectRowAffected: "INSERT 0",
		IsExpectRsNil:         false,
	},
}

func TestSignUpAuthService(t *testing.T) {
	fmt.Println("--------- Authentication Test in process ---------")
	for _, ele := range signUpTestCase {
		t.Run(
			ele.name,
			func(t *testing.T) {
				mockAuthDB := new(mockUtil.MockDB)
				mockTx := new(mockUtil.MockTx)
				reqUsername := "username1"
				reqNickname := "nickname1"
				reqPassword := "password1"
				mockRow := mockUtil.NewMockQueryRow(ele.rowContext)
				mockRow.On("Scan", mock.Anything).Return(ele.scan)
				mockTx.On(
					"QueryRow",
					mock.Anything,
					"INSERT INTO user_setting (darkmode, sound, colorpalettes, font, language) VALUES (0, 0, 0, 1, 1) RETURNING settingid",
					[]interface{}(nil),
				).Return(mockRow)
				mockTx.On("Rollback", mock.Anything).Return(nil)
				mockTx.On("Commit", mock.Anything).Return(nil)
				mockTx.On(
					"Exec",
					mock.Anything,
					"INSERT INTO userauth ( username , nickname , password , setting_id ) VALUES ( $1 , $2 , $3 , $4 )",
					[]interface{}{reqUsername, reqNickname, reqPassword, "2"},
				).Return(pgconn.NewCommandTag(ele.ExecExpectRowAffected), nil)
				authSvc := service.NewAuthService(mockAuthDB)
				userSvc := service.NewUserService(mockAuthDB)

				testRs := authSvc.SignUp(mockTx, userSvc, reqUsername, reqNickname, reqPassword)
				if ele.IsExpectRsNil {
					require.Nil(t, testRs)
				} else {
					require.NotNil(t, testRs)
				}
			},
		)
	}
}

func TestSigningTokenService(t *testing.T) {
	t.Run(
		"Test Signing Token with: Provide all require | Expected No Error",
		func(t *testing.T) {
			mockDB := new(mockUtil.MockDB)
			authSvc := service.NewAuthService(mockDB)
			providedStruct := model.UserData{
				Id:             1,
				Username:       "username1",
				Password:       "password1",
				Nickname:       "nickname1",
				RegisteredDate: time.Now(),
				SettingId:      1,
			}
			testRs, testErrRs := authSvc.SigningToken(providedStruct)
			require.Nil(t, testErrRs)
			require.NotEmpty(t, testRs)
		},
	)
}

func TestParseJwt(t *testing.T) {
	t.Run(
		"Test ParseJwt Service with: Use SignToken Svc For Parsing | Expect No Error",
		func(t *testing.T) {

			mockDB := new(mockUtil.MockDB)
			authSvc := service.NewAuthService(mockDB)
			providedStruct := model.UserData{
				Id:             1,
				Username:       "username1",
				Password:       "password1",
				Nickname:       "nickname1",
				RegisteredDate: time.Now(),
				SettingId:      1,
			}
			testJwtSignedRs, testErrRs := authSvc.SigningToken(providedStruct)
			require.Nil(t, testErrRs)
			newCookie := &http.Cookie{
				Name:     "auth_token",
				Value:    testJwtSignedRs,
				HttpOnly: true,
				Secure:   false,
				Path:     "/",
				SameSite: http.SameSiteLaxMode,
			}

			testJwtClaim, testParseErr := authSvc.ParseJwt(newCookie)
			require.Nil(t, testParseErr)
			require.NotNil(t, testJwtClaim)
		},
	)
}
