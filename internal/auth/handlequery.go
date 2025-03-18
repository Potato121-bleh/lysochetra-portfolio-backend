package auth

import (
	// "backend/internal/domain/service"
	"backend/internal/domain/model"
	"backend/internal/domain/service"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/golang-jwt/jwt/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

var JwtClaimsContextKey ContextKey = "jwtToken"

func newHandler() {

}

type HandleQuery struct {
	DB         *pgxpool.Pool
	UserSvc    *service.UserService[model.UserData]
	SettingSvc *service.UserService[model.SettingStruct]
}

type ReqStruct struct {
	SettingId int `json:"settingId"`
}

// type SettingStruct struct {
// 	SettingId     int `json:"settingid"`
// 	Darkmode      int `json:"darkmode"`
// 	Sound         int `json:"sound"`
// 	Colorpalettes int `json:"colorpalettes"`
// 	Font          int `json:"font"`
// 	Language      int `json:"language"`
// }

func (q *HandleQuery) HandleQuerySetting(w http.ResponseWriter, r *http.Request) {
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

	// s.UService.Select()

	// row := s.DB.QueryRow(context.Background(), "SELECT * FROM user_setting WHERE settingid = $1", reqBody.SettingId)

	// var settingidTem int
	// var darkmodeTem int
	// var soundTem int
	// var colorpalettesTem int
	// var fontTem int
	// var languageTem int

	// scanErr := row.Scan(&settingidTem, &darkmodeTem, &soundTem, &colorpalettesTem, &fontTem, &languageTem)
	// if scanErr != nil {
	// 	http.Error(w, "failed to retrieve data: "+scanErr.Error(), http.StatusInternalServerError)
	// 	return
	// }

	// responseStruct := SettingStruct{
	// 	SettingId:     settingidTem,
	// 	Darkmode:      darkmodeTem,
	// 	Sound:         soundTem,
	// 	Colorpalettes: colorpalettesTem,
	// 	Font:          fontTem,
	// 	Language:      languageTem,
	// }

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
func (q *HandleQuery) HandleUpdateSetting(w http.ResponseWriter, r *http.Request) {
	jwtClaims := r.Context().Value(JwtClaimsContextKey).(jwt.MapClaims)
	settingId := int(jwtClaims["SettingId"].(float64))

	fmt.Println("HERE THE SETTING ID: " + string(int(settingId)))

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

	colVal := []int{reqSettingUpdate.Darkmode, reqSettingUpdate.Sound, reqSettingUpdate.Colorpalettes, reqSettingUpdate.Font, reqSettingUpdate.Language}
	// convert into []string{}
	colValFilter := make([]string, len(colVal))
	for i := range colValFilter {
		colValFilter[i] = strconv.Itoa(colVal[i])
	}
	fmt.Println(colValFilter)

	updateErr := q.SettingSvc.Update(
		tx, "user_setting",
		[]string{"darkmode", "Sound", "colorpalettes", "font", "language"},
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

	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("Update Successfully"))
	commitErr := tx.Commit(context.Background())
	if commitErr != nil {
		http.Error(w, "failed to commit setting update.", http.StatusInternalServerError)
		return
	}
}
