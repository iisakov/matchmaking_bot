package config

import (
	"encoding/json"
	"fmt"
	"matchmaking_bot/model"
	"os"
)

var TOKEN string
var PUBLIC_BOT_CHAT int64
var MODERATOR_BOT_CHAT int64

var ADMINS, MODERATORS, CUSTOMERS model.Users

var PAIRS model.Pairs

func CreateBackup(fi string) {
	json, err := json.Marshal(CUSTOMERS)
	if err != nil {
		fmt.Println(err.Error())
	}
	os.WriteFile(fi, json, 0666)
}

func ReadBackup(fo string) {
	f, err := os.ReadFile(fo)
	if err != nil {
		panic(err)
	}
	err = json.Unmarshal(f, &CUSTOMERS)
	if err != nil {
		panic(err)
	}
}
