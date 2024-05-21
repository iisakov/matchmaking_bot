package main

import (
	"log"
	"matchmaking_bot/config"
	"matchmaking_bot/model"
	"os"
	"strconv"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/joho/godotenv"
)

func init() {
	if err := godotenv.Load(".env"); err != nil {
		log.Fatal("Не найден файл .env")
	}
	config.TOKEN = os.Getenv("MATCHMAKER_BOT_TOKEN")
	config.TEST_BOT_CHAT, _ = strconv.Atoi(os.Getenv("TEST_BOT_CHAT"))
	config.MODERATOR_BOT_CHAT, _ = strconv.Atoi(os.Getenv("MODERATOR_BOT_CHAT"))
}

func main() {

	bot, err := tgbotapi.NewBotAPI(config.TOKEN)
	if err != nil {
		log.Panic(err)
	}
	myBot := model.NewTgBot(bot)

	myBot.Bot.Debug = true

	// myBot.SendMsgById(int64(config.TEST_BOT_CHAT), "Копия бота matchmaker_by_artisan_bot запущена.", myBot.Stage.StageName)
	msg := tgbotapi.NewMessage(int64(config.TEST_BOT_CHAT), "Копия бота matchmaker_by_artisan_bot запущена.")
	myBot.Bot.Send(msg)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := myBot.Bot.GetUpdatesChan(u)

	for update := range updates {
		// Приватные сообщения от пользователей
		if update.Message != nil {
			switch myBot.Stage.StageType {
			case 1: // Регистрация
				if !config.CUSTOMERS.IsExistUserById(update.Message.From.ID) {
					config.CUSTOMERS = append(config.CUSTOMERS, model.NewUser(update))
				} else {
					config.CUSTOMERS.FindUserByIdAndUpdateAlias(update.Message.From.ID, update.Message.Text)
				}
			default:
				myBot.SendMsgById(update.Message.From.ID, "Проше прощения, пока мне нечего на это ответить.")
				myBot.SendMsgById(int64(config.TEST_BOT_CHAT), "Проше прощения, пока мне нечего на это ответить.")
			}

		}

		// Сообщения от Модераторов и Админов
		if update.ChannelPost != nil {
			if !update.ChannelPost.IsCommand() {
				if update.ChannelPost.SenderChat.ID == int64(config.MODERATOR_BOT_CHAT) {
					for _, user := range config.CUSTOMERS {
						myBot.SendMsgById(int64(user.UserChat_id), "Команда [by_artisan]:", update.ChannelPost.Text)
						myBot.SendMsgById(int64(config.MODERATOR_BOT_CHAT), "Пользователю: "+user.UserLogin+", он же "+user.UserAlias, "Отправлено сообщение")
						myBot.SendMsgById(int64(config.TEST_BOT_CHAT), "Пользователю: "+user.UserLogin+", он же "+user.UserAlias, "Отправлено сообщение")
					}
				}
			} else {
				switch update.ChannelPost.Command() {
				case "help":
					myBot.SendMsgById(int64(config.MODERATOR_BOT_CHAT), "Здеь будет подсказка для модераторов")
					myBot.SendMsgById(int64(config.TEST_BOT_CHAT), "case \"help\"")
					continue

				case "stage_up":
					myBot.Stage, err = myBot.Stage.Up()
					if err != nil {
						myBot.SendMsgById(int64(config.MODERATOR_BOT_CHAT), err.Error(), myBot.Stage.StageName)
						myBot.SendMsgById(int64(config.TEST_BOT_CHAT), "err", err.Error())
						continue
					}
					myBot.SendMsgById(int64(config.MODERATOR_BOT_CHAT), "Мы перешли на следущий этап.", myBot.Stage.StageName)
					myBot.SendMsgById(int64(config.TEST_BOT_CHAT), "Мы перешли на следущий этап.", myBot.Stage.StageName)
					continue

				case "stage_down":
					myBot.Stage, err = myBot.Stage.Down()
					if err != nil {
						myBot.SendMsgById(int64(config.MODERATOR_BOT_CHAT), err.Error(), myBot.Stage.StageName)
						myBot.SendMsgById(int64(config.TEST_BOT_CHAT), "err", err.Error())
						continue
					}
					myBot.SendMsgById(int64(config.MODERATOR_BOT_CHAT), "Мы вернулись на предыдущий этап.", myBot.Stage.StageName)
					myBot.SendMsgById(int64(config.TEST_BOT_CHAT), "Мы вернулись на предыдущий этап.", myBot.Stage.StageName)
					continue
				}

			}

		}

	}
}
