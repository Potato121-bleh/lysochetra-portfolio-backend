package middlewareTest

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"profile-portfolio/internal/application"
	"profile-portfolio/internal/domain/model"
	"profile-portfolio/internal/middleware"
	mockUtil "profile-portfolio/test/mock"
	testUtilTool "profile-portfolio/test/util"
	"testing"
	"time"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestMiddlewareCORSValidation(t *testing.T) {

	req := httptest.NewRequest("OPTIONS", "/", nil)
	rec := httptest.NewRecorder()
	handler := middleware.MiddlewareCORSValidate(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
		},
		))
	handler(rec, req)
	t.Run(
		"Test CORS validation, Expect No Error",
		func(t *testing.T) {
			require.NotNil(t, rec.Header().Get("Access-Control-Allow-Origin"))
			require.NotNil(t, rec.Header().Get("Access-Control-Allow-Methods"))
			require.NotNil(t, rec.Header().Get("Access-Control-Allow-Headers"))
			require.NotNil(t, rec.Header().Get("Access-Control-Allow-Credentials"))

		},
	)

}

func TestMiddlewareValidateAuth(t *testing.T) {
	reqData := map[string]string{"username": "username1", "password": "password1"}
	reqJson, marshalReqErr := json.Marshal(reqData)
	require.Nil(t, marshalReqErr)

	req := httptest.NewRequest("POST", "/", bytes.NewBuffer(reqJson))
	rec := httptest.NewRecorder()
	mockDB := new(mockUtil.MockDB)
	mockTx := new(mockUtil.MockTx)
	mockRow := mockUtil.NewMockRow([][]interface{}{{
		1,
		"username1",
		"nickname1",
		"password1",
		time.Now(),
		1,
	}})
	mockScanArg := testUtilTool.CountMockAnything(model.UserData{})
	mockRow.On("Scan", mockScanArg...).Return(nil)
	mockTx.On("Rollback", mock.Anything).Return(nil)
	mockTx.On("Commit", mock.Anything).Return(nil)
	mockTx.On("Query", mock.Anything, "SELECT * FROM userauth WHERE LOWER(username) = $1", []interface{}{"username1"}).Return(mockRow, nil)
	mockDB.On("Begin", mock.Anything).Return(mockTx, nil)
	testHandler := middleware.MiddlewareValidateAuth(
		http.HandlerFunc(
			func(w http.ResponseWriter, r *http.Request) {
				t.Run(
					"Test Auth Middleware Validation, Expect No Error",
					func(t *testing.T) {
						reqContext := r.Context().Value(application.NewUserKey)
						require.NotNil(t, reqContext)
					},
				)
				w.WriteHeader(http.StatusOK)
			}),
		mockDB,
	)

	testHandler(rec, req)

}

func TestMiddlewareCSRFCheck(t *testing.T) {
	handler := middleware.MiddlewareCSRFCheck(
		http.HandlerFunc(
			func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
			},
		),
	)
	req := httptest.NewRequest("GET", "/", nil)
	rec := httptest.NewRecorder()
	handler(rec, req)
	require.Equal(t, http.StatusUnauthorized, rec.Result().StatusCode)
}
