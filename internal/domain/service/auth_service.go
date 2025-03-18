package service

import (
	"backend/internal/domain/model"
	"context"
	"fmt"
	"strconv"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type AuthService struct {
	db         *pgxpool.Pool
	userSvc    *UserService[model.UserData]
	settingSvc *UserService[model.SettingStruct]
}

func (s *AuthService) SignUp(tx pgx.Tx, reqUsername string, reqNickname string, reqPassword string) error {

	// we insert setting, and query id from setting
	newSettingRow := tx.QueryRow(context.Background(), "INSERT INTO user_setting (darkmode, sound, colorpalettes, font, language) VALUES (0, 0, 0, 1, 1)	RETURNING settingid")
	var latestSettingId int
	queriedSettingIdErr := newSettingRow.Scan(&latestSettingId)
	if queriedSettingIdErr != nil {
		return fmt.Errorf("failed to query settingid")
	}

	// we insert new user with setting id we currently
	fmt.Println([]string{reqUsername, reqNickname, reqPassword, strconv.Itoa(latestSettingId)})
	// insertUserErr := s.userSvc.SqlInsert(
	// 	tx,
	// 	"userauth",
	// 	[]string{"username", "nickname", "password", "setting_id"},
	// 	[]string{reqUsername, reqNickname, reqPassword, strconv.Itoa(latestSettingId)},
	// )
	insertUserErr := s.userSvc.Insert(
		tx,
		"userauth",
		[]string{"username", "nickname", "password", "setting_id"},
		[]string{reqUsername, reqNickname, reqPassword, strconv.Itoa(latestSettingId)},
	)
	if insertUserErr != nil {
		return fmt.Errorf("failed to insert userauth (%v)", insertUserErr.Error())
	}

	return nil

}
