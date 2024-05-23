package config

import (
	"matchmaking_bot/model"
	"os"
	"strconv"
)

var TOKEN string = os.Getenv("MATCHMAKER_BOT_TOKEN")
var PUBLIC_BOT_CHAT, _ = strconv.Atoi(os.Getenv("PUBLIC_BOT_CHAT"))
var MODERATOR_BOT_CHAT, _ = strconv.Atoi(os.Getenv("MODERATOR_BOT_CHAT"))

var ADMINS, MODERATORS, CUSTOMERS model.Users

var PAIRS model.Pairs
