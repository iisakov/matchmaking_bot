package main

import (
	"log"
	"matchmaking_bot/config"
	"matchmaking_bot/model"
	"matchmaking_bot/stl/tgstl"
	"os"
	"strconv"

	tgbotapi "github.com/iisakov/telegram-bot-api"
	"github.com/joho/godotenv"
)

func init() {
	if err := godotenv.Load(".env"); err != nil {
		log.Fatal("Не найден файл .env")
	}
	config.TOKEN = os.Getenv("MATCHMAKER_BOT_TOKEN")
	config.PUBLIC_BOT_CHAT, _ = strconv.Atoi(os.Getenv("PUBLIC_BOT_CHAT"))
	config.MODERATOR_BOT_CHAT, _ = strconv.Atoi(os.Getenv("MODERATOR_BOT_CHAT"))
}

func main() {

	bot, err := tgbotapi.NewBotAPI(config.TOKEN)
	if err != nil {
		log.Panic(err)
	}
	myBot := model.NewTgBot(bot)

	myBot.Bot.Debug = false

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := myBot.Bot.GetUpdatesChan(u)

	for update := range updates {
		// Приватные сообщения от пользователей
		if update.Message != nil {
			if !update.Message.IsCommand() {
				tgstl.HandleMessagesText(*update.Message, myBot)
			} else {
				tgstl.HandleMessageCommands(*update.Message, myBot)
			}
		}

		// Сообщения от Модераторов и Админов
		if update.ChannelPost != nil {
			if !update.ChannelPost.IsCommand() {
				tgstl.HandleChannelPostText(*update.ChannelPost, myBot)
			} else {
				myBot = tgstl.HandleChannelPostCommands(*update.ChannelPost, myBot)
			}
		}

		if update.CallbackQuery != nil {
			tgstl.HandleCallbackQuery(*update.CallbackQuery, myBot)
		}
	}
}
