package handlerTest

import (
	"bytes"
	"encoding/json"
	"net/http/httptest"
	"profile-portfolio/internal/application"
	"profile-portfolio/internal/domain/model"
	"profile-portfolio/internal/domain/service"
	mockUtil "profile-portfolio/test/mock"
	testUtilTool "profile-portfolio/test/util"
	"testing"
	"time"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestSettingHandlerQuery(t *testing.T) {
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
