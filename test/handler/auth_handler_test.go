package handlerTest

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"profile-portfolio/internal/application"
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
	name        string
	isFoundUser error
}

func TestSigningTokenHandler(t *testing.T) {
	req := httptest.NewRequest("POST", "/", nil)
	rec := httptest.NewRecorder()
	mockDB := new(mockUtil.MockDB)
	authSvc := service.NewAuthService(mockDB)
	prepModel := model.UserData{
		Id:             1,
		Username:       "username1",
		Nickname:       "nickname1",
		Password:       "password1",
		RegisteredDate: time.Now(),
		SettingId:      1,
	}
	prepContext := context.WithValue(req.Context(), application.NewUserKey, prepModel)
	req = req.WithContext(prepContext)
	authHandler := application.AuthHandler{
		AuthSvc: authSvc,
	}
	authHandler.HandleSigningToken(rec, req)
	resp := rec.Result()
	defer resp.Body.Close()
	testRs := resp.Cookies()[0].Value

	t.Run(
		"Test Signing Token Handler: Expected No Error",
		func(t *testing.T) {
			require.NotNil(t, testRs)
		},
	)
}

func TestVerifyTokenHandler(t *testing.T) {
	req := httptest.NewRequest("GET", "/", nil)
	rec := httptest.NewRecorder()
	mockDB := new(mockUtil.MockDB)
	authSvc := service.NewAuthService(mockDB)
	authhandler := application.AuthHandler{
		AuthSvc: authSvc,
	}
	prepModel := model.UserData{
		Id:             1,
		Username:       "username1",
		Nickname:       "nickname1",
		Password:       "password1",
		RegisteredDate: time.Now(),
		SettingId:      1,
	}

	jwtSignature, signErr := authSvc.SigningToken(prepModel)
	require.Nil(t, signErr)
	newCookie := &http.Cookie{
		Name:     "auth_token",
		Value:    jwtSignature,
		HttpOnly: true,
		Secure:   false,
		Path:     "/",
		SameSite: http.SameSiteLaxMode,
	}
	req.AddCookie(newCookie)
	authhandler.HandleVerifyToken(rec, req)

	t.Run(
		"Test Verify Token Handler: Expect No Error",
		func(t *testing.T) {
			resp := rec.Result()
			defer resp.Body.Close()
			var respData model.UserData
			decodeTestRsErr := json.NewDecoder(resp.Body).Decode(&respData)
			require.Nil(t, decodeTestRsErr)
			require.Equal(t, "username1", respData.Username)
		},
	)

}

func TestRetrieveCSRFKeyHandler(t *testing.T) {
	req := httptest.NewRequest("GET", "/", nil)
	rec := httptest.NewRecorder()
	mockDB := new(mockUtil.MockDB)
	authSvc := service.NewAuthService(mockDB)
	prepModel := model.UserData{
		Id:             1,
		Username:       "username1",
		Nickname:       "nickname1",
		Password:       "password1",
		RegisteredDate: time.Now(),
		SettingId:      1,
	}
	jwtSignature, signErr := authSvc.SigningToken(prepModel)
	require.Nil(t, signErr)
	newCookie := &http.Cookie{
		Name:     "auth_token",
		Value:    jwtSignature,
		HttpOnly: true,
		Secure:   false,
		Path:     "/",
		SameSite: http.SameSiteLaxMode,
	}
	req.AddCookie(newCookie)
	prepContext := context.WithValue(req.Context(), application.NewUserKey, prepModel)
	req = req.WithContext(prepContext)

	authHandler := application.AuthHandler{
		AuthSvc: authSvc,
	}
	authHandler.RetrieveCSRFKey(rec, req)
	t.Run(
		"Test RetrieveCSRF Key Handler: Expect No Error",
		func(t *testing.T) {
			body, readErr := io.ReadAll(rec.Body)
			require.Nil(t, readErr)

			require.NotNil(t, string(body))
		},
	)

}

func TestSignupHandler(t *testing.T) {

	testCase := []signUpTestCaseType{
		{
			name:        "Test Sign up Handler: Expected No Error",
			isFoundUser: fmt.Errorf("No row found"),
		},
		{
			name:        "Test Sign up Handler: Expected Error",
			isFoundUser: nil,
		},
	}

	prepReqData := model.UserData{
		Username: "username1",
		Nickname: "nickname1",
		Password: "password1",
	}
	reqDataJson, marshalErr := json.Marshal(prepReqData)
	require.Nil(t, marshalErr)

	for _, ele := range testCase {
		t.Run(
			ele.name,
			func(t *testing.T) {
				mockDB := new(mockUtil.MockDB)
				mockTx := new(mockUtil.MockTx)
				authSvc := service.NewAuthService(mockDB)
				userSvc := service.NewUserService(mockDB)
				mockDBRow := mockUtil.NewMockQueryRow([]interface{}{})
				mockTxRow := mockUtil.NewMockQueryRow([]interface{}{
					2,
				})

				mockDBRow.On("Scan", mock.Anything).Return(ele.isFoundUser)
				mockDB.On("Begin", mock.Anything).Return(mockTx, nil)
				mockDB.On(
					"QueryRow",
					mock.Anything,
					"SELECT userid FROM userauth WHERE username = $1",
					[]interface{}{prepReqData.Username},
				).Return(mockDBRow)

				mockTx.On(
					"QueryRow",
					mock.Anything,
					"INSERT INTO user_setting (darkmode, sound, colorpalettes, font, language) VALUES (0, 0, 0, 1, 1) RETURNING settingid",
					[]interface{}(nil),
				).Return(mockTxRow)
				mockTx.On(
					"Exec",
					mock.Anything,
					"INSERT INTO userauth ( username , nickname , password , setting_id ) VALUES ( $1 , $2 , $3 , $4 )",
					[]interface{}{prepReqData.Username, prepReqData.Nickname, prepReqData.Password, "2"},
				).Return(pgconn.NewCommandTag("INSERT 1"), nil)
				mockTx.On("Commit", mock.Anything).Return(nil)
				mockTx.On("Rollback", mock.Anything).Return(nil)

				mockTxRow.On("Scan", mock.Anything).Return(nil)

				req := httptest.NewRequest("POST", "/", bytes.NewBuffer(reqDataJson))
				rec := httptest.NewRecorder()

				authHandler := application.AuthHandler{
					DB:      mockDB,
					AuthSvc: authSvc,
					UserSvc: userSvc,
				}
				authHandler.HandleSignup(rec, req)
				if ele.isFoundUser != nil {
					require.Equal(t, http.StatusCreated, rec.Code)
				} else {
					require.Equal(t, http.StatusBadRequest, rec.Code)
				}
			},
		)
	}
}

func TestLogoutHandler(t *testing.T) {
	req := httptest.NewRequest("GET", "/", nil)
	rec := httptest.NewRecorder()

	authHandler := application.AuthHandler{}
	authHandler.HandleLogout(rec, req)
	resCookies := rec.Result().Cookies()

	t.Run(
		"Test Log out Service: No Error Expected",
		func(t *testing.T) {
			require.Equal(t, http.StatusOK, rec.Code)
			require.Equal(t, len(resCookies), 1)
		},
	)
}
