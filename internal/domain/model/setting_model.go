package model

type SettingStruct struct {
	SettingId     int `json:"settingid"`
	Darkmode      int `json:"darkmode"`
	Sound         int `json:"sound"`
	Colorpalettes int `json:"colorpalettes"`
	Font          int `json:"font"`
	Language      int `json:"language"`
}
