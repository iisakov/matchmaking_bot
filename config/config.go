package config

import (
	"encoding/json"
	"errors"
	"fmt"
	"matchmaking_bot/model"
	"os"
)

var TOKEN string
var PUBLIC_BOT_CHAT int64
var MODERATOR_BOT_CHAT int64

var ADMINS, MODERATORS, CUSTOMERS model.Users

var PAIRS model.Pairs

type Backupable interface {
	CreateBackup()
	ReadBackup()
}

func CreateBackup() {
	var l = map[string]Backupable{
		"ADMINS":     ADMINS,
		"MODERATORS": MODERATORS,
		"CUSTOMERS":  CUSTOMERS,
		"PAIRS":      PAIRS}
	for k, v := range l {
		json, err := json.Marshal(v)
		if err != nil {
			fmt.Println(err.Error())
		}
		os.WriteFile(fmt.Sprintf("backup%s.json", k), json, 0666)
	}
}

func ReadBackup() {
	var l = map[string]Backupable{
		"ADMINS":     &ADMINS,
		"MODERATORS": &MODERATORS,
		"CUSTOMERS":  &CUSTOMERS,
		"PAIRS":      &PAIRS}
	for k, v := range l {

		if _, err := os.Stat(fmt.Sprintf("backup%s.json", k)); errors.Is(err, os.ErrNotExist) {
			fmt.Println(err)
		} else {
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
}
