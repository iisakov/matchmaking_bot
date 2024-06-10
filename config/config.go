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

func CreateBackup() {
	var l = map[string]model.Users{"ADMINS": ADMINS, "MODERATORS": MODERATORS, "CUSTOMERS": CUSTOMERS}
	for k, v := range l {
		json, err := json.Marshal(v)
		if err != nil {
			fmt.Println(err.Error())
		}
		os.WriteFile(fmt.Sprintf("backup%s.json", k), json, 0666)
	}
}

func ReadBackup() {
	var l = map[string]*model.Users{"ADMINS": &ADMINS, "MODERATORS": &MODERATORS, "CUSTOMERS": &CUSTOMERS}
	for k, v := range l {
		f, err := os.ReadFile(fmt.Sprintf("backup%s.json", k))
		if err != nil {
			panic(err)
		}
		err = json.Unmarshal(f, v)
		if err != nil {
			panic(err)
		}
		l[k] = v
	}
}
