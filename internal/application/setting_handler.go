package application

import (
	// "profile-portfolio/internal/domain/service"
	"context"
	"encoding/json"
	"net/http"
	"profile-portfolio/internal/db"
	"profile-portfolio/internal/domain/model"
	"profile-portfolio/internal/domain/service"
	"strconv"

	"github.com/golang-jwt/jwt/v5"
)

var JwtClaimsContextKey ContextKey = "jwtToken"

type SettingHandler struct {
	DB         db.Database
	UserSvc    service.UserServiceI
	SettingSvc service.SettingServiceI
}

func (q *SettingHandler) HandleQuerySetting(w http.ResponseWriter, r *http.Request) {
	var reqBody = model.UserData{}
	reqBodyErr := json.NewDecoder(r.Body).Decode(&reqBody)
	if reqBodyErr != nil {
		http.Error(w, "failed to decode request: "+reqBodyErr.Error(), http.StatusBadRequest)
		return
	}

	if reqBody.SettingId == 0 {
		http.Error(w, "please include settingId for retrieving data", http.StatusBadRequest)
		return
	}

	responseStruct, queryErr := q.SettingSvc.Select(nil, "user_setting", "settingid", strconv.Itoa(reqBody.SettingId))
	if queryErr != nil {
		http.Error(w, "failed to perform transaction", http.StatusInternalServerError)
		return
	}

	encodeRespErr := json.NewEncoder(w).Encode(responseStruct[0])
	if encodeRespErr != nil {
		http.Error(w, "failed to encode the response data", http.StatusInternalServerError)
		return
	}
}

// This update func will work when you request with all data such as every field of setting except settingId
func (q *SettingHandler) HandleUpdateSetting(w http.ResponseWriter, r *http.Request) {
	jwtClaims := r.Context().Value(JwtClaimsContextKey).(jwt.MapClaims)
	settingId := int(jwtClaims["SettingId"].(float64))
	reqSettingUpdate := model.SettingStruct{}
	decodeReqBodyErr := json.NewDecoder(r.Body).Decode(&reqSettingUpdate)
	if decodeReqBodyErr != nil {
		http.Error(w, "failed to decode request body", http.StatusBadRequest)
		return
	}

	tx, txErr := q.DB.Begin(context.Background())
	if txErr != nil {
		http.Error(w, "failed to update setting", http.StatusInternalServerError)
		return
	}

	colVal := []int{
		reqSettingUpdate.Darkmode,
		reqSettingUpdate.Sound,
		reqSettingUpdate.Colorpalettes,
		reqSettingUpdate.Font,
		reqSettingUpdate.Language,
	}

	colValFilter := make([]string, len(colVal))
	for i := range colValFilter {
		colValFilter[i] = strconv.Itoa(colVal[i])
	}

	updateErr := q.SettingSvc.Update(
		tx, "user_setting",
		[]string{"darkmode", "sound", "colorpalettes", "font", "language"},
		colValFilter,
		"settingid",
		strconv.Itoa(settingId),
	)

	if updateErr != nil {
		rollbackErr := tx.Rollback(context.Background())
		if rollbackErr != nil {
			http.Error(w, "failed to update setting", http.StatusInternalServerError)
			return
		}

		http.Error(w, "failed to update setting", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Update Successfully"))
	commitErr := tx.Commit(context.Background())
	if commitErr != nil {
		http.Error(w, "failed to commit setting update.", http.StatusInternalServerError)
		return
	}
}
