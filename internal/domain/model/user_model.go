package model

import "time"

type UserData struct {
	Id             int       `json:"userId"`
	Username       string    `json:"userName"`
	Nickname       string    `json:"nickname"`
	Password       string    `json:"password"`
	RegisteredDate time.Time `json:"registered_date"`
	SettingId      int       `json:"settingId"`
}

type SettingStruct struct {
	SettingId     int `json:"settingid"`
	Darkmode      int `json:"darkmode"`
	Sound         int `json:"sound"`
	Colorpalettes int `json:"colorpalettes"`
	Font          int `json:"font"`
	Language      int `json:"language"`
}
