package main

import (
	"fmt"
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
	config.CUSTOMERS = config.MockUsers(config.CUSTOMERS) // моковые пользователи для проверки

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
					myBot.SendMsgById(update.Message.From.ID, "Прошу прощения, пока мне нечего на это ответить.")
				}
			} else {
				switch update.Message.Command() {
				case "start":
					if myBot.Stage.StageType == 1 {
						myBot.SendMsgById(update.Message.From.ID, "Привет, настало время придумать себе псевдоним.")
						continue
					}
					myBot.SendMsgById(update.Message.From.ID, "Прошу прощения, пока мне нечего на это ответить.")
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
					fallthrough
				case "h":
					myBot.SendMsgById(int64(config.MODERATOR_BOT_CHAT), "Здеь будет подсказка для модераторов")
					continue

				case "chack_stages_text":
					fallthrough
				case "cst":
					myBot.SendMsgById(int64(config.MODERATOR_BOT_CHAT), myBot.Stage.StageName, myBot.Stage.StagesText)
					continue

				case "send_stages_text":
					fallthrough
				case "sst":
					myBot.SendMsgById(int64(config.MODERATOR_BOT_CHAT), myBot.Stage.StageName, myBot.Stage.StagesText)
					myBot.SendMsgById(int64(config.TEST_BOT_CHAT), myBot.Stage.StageName, myBot.Stage.StagesText)
					continue

				case "stage_up":
					fallthrough
				case "su":
					myBot.Stage, err = myBot.Stage.Up()
					if err != nil {
						myBot.SendMsgById(int64(config.MODERATOR_BOT_CHAT), err.Error(), myBot.Stage.StageName, strconv.Itoa(myBot.Stage.StageType))
						continue
					}
					myBot.SendMsgById(int64(config.MODERATOR_BOT_CHAT), "Мы перешли на следущий этап.", myBot.Stage.StageName, strconv.Itoa(myBot.Stage.StageType))
					continue

				case "stage_down":
					fallthrough
				case "sd":
					myBot.Stage, err = myBot.Stage.Down()
					if err != nil {
						myBot.SendMsgById(int64(config.MODERATOR_BOT_CHAT), err.Error(), myBot.Stage.StageName, strconv.Itoa(myBot.Stage.StageType))
						continue
					}
					myBot.SendMsgById(int64(config.MODERATOR_BOT_CHAT), "Мы вернулись на предыдущий этап.", myBot.Stage.StageName, strconv.Itoa(myBot.Stage.StageType))
					continue

				case "question_next":
					fallthrough
				case "qn":

					model.BotQuestions, err = model.BotQuestions.Next()
					if err != nil {
						myBot.SendMsgById(int64(config.MODERATOR_BOT_CHAT), err.Error(), model.BotQuestions.GetCurentQuestion().Text)
						continue
					}
					myBot.SendMsgById(int64(config.MODERATOR_BOT_CHAT), "Следующий вопрос звучит так.", model.BotQuestions.GetCurentQuestion().Text)
					continue

				case "question_back":
					fallthrough
				case "qb":
					model.BotQuestions, err = model.BotQuestions.Back()
					if err != nil {
						myBot.SendMsgById(int64(config.MODERATOR_BOT_CHAT), err.Error(), model.BotQuestions.GetCurentQuestion().Text)
						continue
					}
					myBot.SendMsgById(int64(config.MODERATOR_BOT_CHAT), "Предыдущий вопрос звучит так.", model.BotQuestions.GetCurentQuestion().Text)
					continue

				case "send_question":
					fallthrough
				case "sq":
					if myBot.Stage.StageType == 2 {
						for _, user := range config.CUSTOMERS {
							myBot.SendMsgWithInleneKeyboardById(
								user.UserChat_id,
								model.BotQuestions.GetCurentQuestion().Markup,
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

				case "create_pair":
					fallthrough
				case "cp":
					config.PAIRS = model.CheckCompatibility(config.CUSTOMERS.GetUsersByGender(1), config.CUSTOMERS.GetUsersByGender(0))
					myBot.SendMsgById(
						int64(config.MODERATOR_BOT_CHAT),
						"Создали пары:",
						config.PAIRS.GetPairs(),
					)
					continue

				case "get_pair":
					fallthrough
				case "gp":
					myBot.SendMsgById(
						int64(config.MODERATOR_BOT_CHAT),
						config.PAIRS.GetPairs(),
					)

					continue

				case "get_users":
					fallthrough
				case "gus":
					if len(config.CUSTOMERS) > 0 {
						myBot.SendMsgById(
							int64(config.MODERATOR_BOT_CHAT),
							config.CUSTOMERS.GetUsersByGender(1).GetUsers(),
							config.CUSTOMERS.GetUsersByGender(0).GetUsers(),
						)
					} else {
						myBot.SendMsgById(
							int64(config.MODERATOR_BOT_CHAT),
							"Нет зарегистрированных пользователей.",
						)
					}

					continue
				}
			}
		}

		if update.CallbackQuery != nil {
			if myBot.Stage.StageType == 2 {
				if config.CUSTOMERS.IsExistUserById(update.CallbackQuery.From.ID) {
					switch {
					case model.BotQuestions.QuestionsList.IsExistQuestionOptionsByName(update.CallbackQuery.Message.Text, "gender"):
						config.CUSTOMERS.FindUserByIdSetGender(update.CallbackQuery.From.ID, update.CallbackQuery.Data)
						myBot.SendMsgById(update.CallbackQuery.From.ID, "Отлично, теперь мы знаем какого ты поля.")
						myBot.SendMsgById(update.CallbackQuery.From.ID, config.CUSTOMERS.FindUserById(update.CallbackQuery.From.ID).Answers...)
						continue

					case model.BotQuestions.QuestionsList.IsExistQuestionOptionsByName(update.CallbackQuery.Message.Text, "onlyOne"):
						fmt.Println(model.BotQuestions.QuestionsList.GetAnswersByQuestionName(update.CallbackQuery.Message.Text))
						config.CUSTOMERS.FindUserByIdAndUpdateAnswer(
							update.CallbackQuery.From.ID,
							model.BotQuestions.QuestionsList.GetAnswersByQuestionName(update.CallbackQuery.Message.Text),
							update.CallbackQuery.Data)
						myBot.SendMsgById(update.CallbackQuery.From.ID, "Отлично, Можешь изменить ответ если хочешь.")
						myBot.SendMsgById(update.CallbackQuery.From.ID, config.CUSTOMERS.FindUserById(update.CallbackQuery.From.ID).Answers...)
						continue

					default:
						config.CUSTOMERS.FindUserByIdAndAddAnswer(update.CallbackQuery.From.ID, update.CallbackQuery.Data)
						myBot.SendMsgById(update.CallbackQuery.From.ID, "Отлично, На этот вопрос можно ответить несколько раз, Выбирай хоть все варианты.")
						myBot.SendMsgById(update.CallbackQuery.From.ID, config.CUSTOMERS.FindUserById(update.CallbackQuery.From.ID).Answers...)
						continue
					}
				}
			}
		}
	}
}
