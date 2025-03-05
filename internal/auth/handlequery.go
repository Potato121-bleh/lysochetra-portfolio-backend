package auth

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/golang-jwt/jwt/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

var JwtClaimsContextKey ContextKey = "jwtToken"

type HandleQuery struct {
	DB *pgxpool.Pool
}

type ReqStruct struct {
	SettingId int `json:"settingId"`
}

type settingStruct struct {
	SettingId     int `json:"settingid"`
	Darkmode      int `json:"darkmode"`
	Sound         int `json:"sound"`
	Colorpalettes int `json:"colorpalettes"`
	Font          int `json:"font"`
	Language      int `json:"language"`
}

func (s *HandleQuery) HandleQuerySetting(w http.ResponseWriter, r *http.Request) {
	var reqBody = ReqStruct{}
	reqBodyErr := json.NewDecoder(r.Body).Decode(&reqBody)
	if reqBodyErr != nil {
		http.Error(w, "failed to decode request: "+reqBodyErr.Error(), http.StatusBadRequest)
		return
	}

	if reqBody.SettingId == 0 {
		http.Error(w, "please include settingId for retrieving data", http.StatusBadRequest)
		return
	}

	row := s.DB.QueryRow(context.Background(), "SELECT * FROM user_setting WHERE settingid = $1", reqBody.SettingId)

	var settingidTem int
	var darkmodeTem int
	var soundTem int
	var colorpalettesTem int
	var fontTem int
	var languageTem int

	scanErr := row.Scan(&settingidTem, &darkmodeTem, &soundTem, &colorpalettesTem, &fontTem, &languageTem)
	if scanErr != nil {
		http.Error(w, "failed to retrieve data: "+scanErr.Error(), http.StatusInternalServerError)
		return
	}

	responseStruct := settingStruct{
		SettingId:     settingidTem,
		Darkmode:      darkmodeTem,
		Sound:         soundTem,
		Colorpalettes: colorpalettesTem,
		Font:          fontTem,
		Language:      languageTem,
	}

	encodeRespErr := json.NewEncoder(w).Encode(responseStruct)
	if encodeRespErr != nil {
		http.Error(w, "failed to encode the response data", http.StatusInternalServerError)
		return
	}
}

// This update func will work when you request with all data such as every field of setting except settingId
func (s *HandleQuery) HandleUpdateSetting(w http.ResponseWriter, r *http.Request) {
	jwtClaims := r.Context().Value(JwtClaimsContextKey).(jwt.MapClaims)
	settingId := int(jwtClaims["SettingId"].(float64))

	fmt.Println("HERE THE SETTING ID: " + string(int(settingId)))

	reqSettingUpdate := settingStruct{}
	decodeReqBodyErr := json.NewDecoder(r.Body).Decode(&reqSettingUpdate)
	if decodeReqBodyErr != nil {
		http.Error(w, "failed to decode request body", http.StatusBadRequest)
		return
	}

	tx, txErr := s.DB.Begin(context.Background())
	if txErr != nil {
		http.Error(w, "failed to update setting", http.StatusInternalServerError)
		return
	}

	//you can just update with all of that direct data in request body because we expect the frontend will does the job for us.
	//So frontend has to send an all field such as darkmode, sound, language ... so that we can just take that and do request directly
	//since our frontend already have an object done already to be sent.

	rows, updateTxErr := tx.Exec(context.Background(), "UPDATE user_setting SET darkmode = $1, Sound = $2, colorpalettes = $3, font = $4, language = $5 where settingid = $6;",
		reqSettingUpdate.Darkmode, reqSettingUpdate.Sound, reqSettingUpdate.Colorpalettes, reqSettingUpdate.Font, reqSettingUpdate.Language, settingId)
	if updateTxErr != nil {
		http.Error(w, "failed to update setting", http.StatusInternalServerError)
		tx.Rollback(context.Background())
		return
	}

	if rows.RowsAffected() != 1 {
		http.Error(w, "failed to update setting", http.StatusInternalServerError)
		tx.Rollback(context.Background())
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("Update Successfully"))
	tx.Commit(context.Background())

}
