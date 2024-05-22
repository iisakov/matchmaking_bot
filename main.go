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

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := myBot.Bot.GetUpdatesChan(u)

	for update := range updates {
		// Приватные сообщения от пользователей
		if update.Message != nil {
			if !update.Message.IsCommand() {
				switch myBot.Stage.StageType {
				case 1: // Регистрация
					if !config.CUSTOMERS.IsExistUserById(update.Message.From.ID) {
						config.CUSTOMERS = append(config.CUSTOMERS, model.NewUser(update))
					} else {
						config.CUSTOMERS.FindUserByIdAndUpdateAlias(update.Message.From.ID, update.Message.Text)
					}
					myBot.SendMsgById(update.Message.From.ID, "Отлично, как бы ты себя не назвал, "+update.Message.Text+", таким тебя будет видеть собеседник. У тебя есть ещё пара минут подумать и изменить псевдоним, просто отправь мне сообщение.")
				default:
					myBot.SendMsgById(update.Message.From.ID, "Проше прощения, пока мне нечего на это ответить.")
				}
			} else {
				switch update.Message.Command() {
				case "start":
					if myBot.Stage.StageType == 1 {
						myBot.SendMsgById(update.Message.From.ID, "Привет, настало время придумать себе псевдоним.")
						continue
					}
					myBot.SendMsgById(update.Message.From.ID, "Проше прощения, пока мне нечего на это ответить.")
				}
			}
		}

		// Сообщения от Модераторов и Админов
		if update.ChannelPost != nil {
			if !update.ChannelPost.IsCommand() {
				if update.ChannelPost.SenderChat.ID == int64(config.MODERATOR_BOT_CHAT) {
					for _, user := range config.CUSTOMERS {
						myBot.SendMsgById(int64(user.UserChat_id), "Команда [by_artisan]:", update.ChannelPost.Text)
						myBot.SendMsgById(int64(config.MODERATOR_BOT_CHAT), "Пользователю: "+user.UserLogin+", он же "+user.UserAlias, "Отправлено сообщение")
					}
				}
			} else {
				switch update.ChannelPost.Command() {
				case "help":
					myBot.SendMsgById(int64(config.MODERATOR_BOT_CHAT), "Здеь будет подсказка для модераторов")
					continue

				case "stage_up":
					myBot.Stage, err = myBot.Stage.Up()
					if err != nil {
						myBot.SendMsgById(int64(config.MODERATOR_BOT_CHAT), err.Error(), myBot.Stage.StageName)
						continue
					}
					myBot.SendMsgById(int64(config.MODERATOR_BOT_CHAT), "Мы перешли на следущий этап.", myBot.Stage.StageName)
					myBot.SendMsgById(int64(config.TEST_BOT_CHAT), "Мы перешли на следущий этап.", myBot.Stage.StageName)
					continue

				case "stage_down":
					myBot.Stage, err = myBot.Stage.Down()
					if err != nil {
						myBot.SendMsgById(int64(config.MODERATOR_BOT_CHAT), err.Error(), myBot.Stage.StageName)
						continue
					}
					myBot.SendMsgById(int64(config.MODERATOR_BOT_CHAT), "Мы вернулись на предыдущий этап.", myBot.Stage.StageName)
					myBot.SendMsgById(int64(config.TEST_BOT_CHAT), "Мы вернулись на предыдущий этап.", myBot.Stage.StageName)
					continue

				case "question_next":
					model.BotQuestions, err = model.BotQuestions.Next()
					if err != nil {
						myBot.SendMsgById(int64(config.MODERATOR_BOT_CHAT), err.Error(), model.BotQuestions.GetCurentQuestion().Text)
						continue
					}
					myBot.SendMsgById(int64(config.MODERATOR_BOT_CHAT), "Следующий вопрос звучит так.", model.BotQuestions.GetCurentQuestion().Text, strconv.FormatBool(model.BotQuestions.GetCurentQuestion().MustMatch))
					continue

				case "question_back":
					model.BotQuestions, err = model.BotQuestions.Back()
					if err != nil {
						myBot.SendMsgById(int64(config.MODERATOR_BOT_CHAT), err.Error(), model.BotQuestions.GetCurentQuestion().Text)
						continue
					}
					myBot.SendMsgById(int64(config.MODERATOR_BOT_CHAT), "Предыдущий вопрос звучит так.", model.BotQuestions.GetCurentQuestion().Text)
					continue

				case "send_question":
					if myBot.Stage.StageType == 2 {
						for _, user := range config.CUSTOMERS {
							myBot.SendMsgWithInleneKeyboardById(
								user.UserChat_id,
								model.BotQuestions.GetCurentQuestion().Markup,
								strconv.Itoa(model.BotQuestions.GetQuestionsCounter()),
								model.BotQuestions.GetCurentQuestion().Text)
							myBot.SendMsgById(int64(config.MODERATOR_BOT_CHAT), "Пользователю: "+user.UserLogin+", он же "+user.UserAlias, "Отправлен вопрос.", model.BotQuestions.GetCurentQuestion().Text)
						}
					} else {
						myBot.SendMsgById(
							int64(config.MODERATOR_BOT_CHAT),
							"Текущий этап - "+myBot.Stage.StageName+" - Отправлять вопросы пользователям опрометчиво.",
							strconv.Itoa(model.BotQuestions.GetQuestionsCounter()),
							model.BotQuestions.GetCurentQuestion().Text)
					}

					continue
				}
			}
		}
	}
}
