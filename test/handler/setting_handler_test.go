package handlerTest

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"profile-portfolio/internal/application"
	"profile-portfolio/internal/domain/model"
	"profile-portfolio/internal/domain/service"
	mockUtil "profile-portfolio/test/mock"
	testUtilTool "profile-portfolio/test/util"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestSettingQueryHandler(t *testing.T) {
	timeNow := time.Now()
	reqData := model.UserData{
		Id:             1,
		Username:       "username1",
		Nickname:       "nickname1",
		Password:       "password1",
		RegisteredDate: timeNow,
		SettingId:      1,
	}
	reqJson, marshalReqErr := json.Marshal(reqData)
	require.Nil(t, marshalReqErr)

	req := httptest.NewRequest("POST", "/", bytes.NewBuffer(reqJson))
	rec := httptest.NewRecorder()

	mockSettingDB := new(mockUtil.MockDB)
	mockTx := new(mockUtil.MockTx)
	expectedRs := [][]interface{}{
		{
			1,
			1,
			0,
			1,
			0,
			0,
		},
	}

	mockRow := mockUtil.NewMockRow(expectedRs)
	mockScanArg := testUtilTool.CountMockAnything(model.SettingStruct{})
	mockRow.On("Scan", mockScanArg...).Return(nil)
	mockTx.On("Commit", mock.Anything).Return(nil)
	mockTx.On("Rollback", mock.Anything).Return(nil)
	mockTx.On("Query", mock.Anything, "SELECT * FROM user_setting WHERE settingid = $1", []interface{}{"1"}).Return(mockRow, nil)
	mockSettingDB.On("Begin", mock.Anything).Return(mockTx, nil)
	settingSvc := service.NewSettingService(mockSettingDB)
	settingHandler := application.SettingHandler{
		SettingSvc: settingSvc,
	}
	settingHandler.HandleQuerySetting(rec, req)

	testResp := rec.Result()
	defer testResp.Body.Close()

	var testRs model.SettingStruct
	decodeErr := json.NewDecoder(testResp.Body).Decode(&testRs)
	require.Nil(t, decodeErr)

	t.Run(
		"Test Query Setting Handler, No Error Expected",
		func(t *testing.T) {
			require.Equal(t, expectedRs[0][1], testRs.Darkmode)
		},
	)

}

func TestSettingUpdateHandler(t *testing.T) {
	mockJwtClaim := jwt.MapClaims{
		"Id":        1,
		"Username":  "username1",
		"Password":  "password1",
		"Nickname":  "nickname1",
		"SettingId": float64(1),
		"exp":       jwt.NewNumericDate(time.Now().Add(1 * time.Hour)),
		"iat":       jwt.NewNumericDate(time.Now()),
	}
	prepReqBody := model.SettingStruct{
		SettingId:     1,
		Darkmode:      1,
		Sound:         0,
		Colorpalettes: 1,
		Font:          0,
		Language:      0,
	}

	reqBodyJson, reqToJsonErr := json.Marshal(prepReqBody)
	require.Nil(t, reqToJsonErr)

	req := httptest.NewRequest("POST", "/", bytes.NewBuffer(reqBodyJson))
	claimsContext := context.WithValue(req.Context(), application.JwtClaimsContextKey, mockJwtClaim)
	fmt.Println(claimsContext)
	req = req.WithContext(claimsContext)

	rec := httptest.NewRecorder()

	mockDB := new(mockUtil.MockDB)
	mockSettingDB := new(mockUtil.MockDB)
	mockTx := new(mockUtil.MockTx)
	mockTx.On("Commit", mock.Anything).Return(nil)
	mockTx.On("Rollback", mock.Anything).Return(nil)
	mockTx.On(
		"Exec",
		mock.Anything,
		"UPDATE user_setting SET darkmode = $1 , sound = $2 , colorpalettes = $3 , font = $4 , language = $5 WHERE settingid = $6",
		[]interface{}{"1", "0", "1", "0", "0", "1"},
	).Return(pgconn.NewCommandTag("UPDATE 1"), nil)
	mockDB.On("Begin", mock.Anything).Return(mockTx, nil)

	settingSvc := service.NewSettingService(mockSettingDB)

	settingHandler := application.SettingHandler{
		DB:         mockDB,
		SettingSvc: settingSvc,
	}
	settingHandler.HandleUpdateSetting(rec, req)

	require.Equal(t, http.StatusOK, rec.Result().StatusCode)

}
