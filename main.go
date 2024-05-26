package main

import (
	"fmt"
	"log"
	"matchmaking_bot/config"
	"matchmaking_bot/model"
	"matchmaking_bot/stl"
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
	// config.CUSTOMERS = mock.MockUsers(config.CUSTOMERS) // моковые пользователи для проверки
	// config.PAIRS = mock.MockPairs(config.PAIRS)         // моковые пары для проверки

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
				switch myBot.Stage.StageType {
				case 1: // Регистрация
					if !config.CUSTOMERS.IsExistUserById(update.Message.From.ID) {
						config.CUSTOMERS = append(config.CUSTOMERS, model.NewUser(update))
					} else {
						config.CUSTOMERS.FindUserByIdAndUpdateAlias(update.Message.From.ID, update.Message.Text)
					}
					myBot.DeleteMessegeByIds(update.Message.From.ID, stl.CreateSliceInt(update.Message.MessageID-100, update.Message.MessageID))
					myBot.SendMsgById(update.Message.From.ID, "Отлично, как бы ты себя не назвал, "+update.Message.Text+", таким тебя будет видеть собеседник. У тебя есть ещё пара минут подумать и изменить псевдоним, просто отправь мне сообщение.")
				case 4: // Общение в парах
					conversationPartnerId, ok := config.PAIRS.GetConversationPartner(update.Message.From.ID)
					if !ok {
						myBot.SendMsgById(int64(config.MODERATOR_BOT_CHAT), fmt.Sprintf("Так вышло, что кому-то не досталась пара: %s", config.CUSTOMERS.FindUserById(update.Message.From.ID).UserLogin))
					}
					myBot.SendMsgById(
						conversationPartnerId,
						"Сообщение от "+config.CUSTOMERS.FindUserById(update.Message.From.ID).UserAlias+":",
						update.Message.Text)
				default:
					myBot.DeleteMessegeByIds(update.Message.From.ID, stl.CreateSliceInt(update.Message.MessageID-100, update.Message.MessageID))
					myBot.SendMsgById(update.Message.From.ID, "Прошу прощения, пока мне нечего на это ответить.")
				}
			} else {
				switch update.Message.Command() {
				case "start":
					if myBot.Stage.StageType == 1 {
						myBot.DeleteMessegeByIds(update.Message.From.ID, stl.CreateSliceInt(update.Message.MessageID-100, update.Message.MessageID))
						myBot.SendMsgById(update.Message.From.ID, "Привет, настало время придумать себе псевдоним.")
						continue
					}
					myBot.DeleteMessegeByIds(update.Message.From.ID, stl.CreateSliceInt(update.Message.MessageID-100, update.Message.MessageID))
					myBot.SendMsgById(update.Message.From.ID, "Прошу прощения, пока мне нечего на это ответить.")
				default:
					myBot.DeleteMessegeByIds(update.Message.From.ID, stl.CreateSliceInt(update.Message.MessageID-100, update.Message.MessageID))
					myBot.SendMsgById(update.Message.From.ID, "Прошу прощения, пока мне нечего на это ответить.")
				}
			}
			config.CUSTOMERS.FindUserByIdAndSetLastMessageId(update.Message.From.ID, update.Message.MessageID)
		}

		// Сообщения от Модераторов и Админов
		if update.ChannelPost != nil {
			if !update.ChannelPost.IsCommand() {
				if update.ChannelPost.SenderChat.ID == int64(config.MODERATOR_BOT_CHAT) {
					myBot.SendMsgById(int64(config.PUBLIC_BOT_CHAT), "Сообщение от команды [by_artisan]:", update.ChannelPost.Text)
					for _, user := range config.CUSTOMERS {
						myBot.SendMsgById(int64(user.UserChat_id), "Сообщение от команды [by_artisan]:", update.ChannelPost.Text)
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
					for _, user := range config.CUSTOMERS {
						myBot.DeleteMessegeByIds(user.UserChat_id, stl.CreateSliceInt(user.LastMessageId-50, user.LastMessageId+50))
						myBot.SendMsgById(user.UserChat_id, myBot.Stage.StageName, myBot.Stage.StagesText)
					}
					myBot.SendMsgById(int64(config.MODERATOR_BOT_CHAT), myBot.Stage.StageName, myBot.Stage.StagesText)
					myBot.SendMsgById(int64(config.PUBLIC_BOT_CHAT), myBot.Stage.StageName, myBot.Stage.StagesText)
					continue

				case "stage_up":
					fallthrough
				case "su":
					if myBot.Stage.StageType+1 == 4 && len(config.PAIRS) < 1 {
						myBot.SendMsgById(
							int64(config.MODERATOR_BOT_CHAT),
							"Нельзя переходить к следующему этапу",
							"Распределённых пар:"+strconv.Itoa(len(config.PAIRS)),
						)
						continue
					}
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
							myBot.DeleteMessegeByIds(user.UserChat_id, stl.CreateSliceInt(user.LastMessageId-50, user.LastMessageId+50))
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
						myBot.SendMsgById(update.CallbackQuery.From.ID, "Отлично, теперь мы знаем какого ты пола.")
						continue

					case model.BotQuestions.QuestionsList.IsExistQuestionOptionsByName(update.CallbackQuery.Message.Text, "onlyOne"):
						fmt.Println(model.BotQuestions.QuestionsList.GetAnswersByQuestionName(update.CallbackQuery.Message.Text))
						config.CUSTOMERS.FindUserByIdAndUpdateAnswer(
							update.CallbackQuery.From.ID,
							model.BotQuestions.QuestionsList.GetAnswersByQuestionName(update.CallbackQuery.Message.Text),
							update.CallbackQuery.Data)
						myBot.SendMsgById(update.CallbackQuery.From.ID, "Отлично, Можешь изменить ответ если хочешь.")
						continue

					default:
						config.CUSTOMERS.FindUserByIdAndAddAnswer(update.CallbackQuery.From.ID, update.CallbackQuery.Data)
						myBot.SendMsgById(update.CallbackQuery.From.ID, "Отлично, На этот вопрос можно ответить несколько раз, Выбирай хоть все варианты.")
						continue
					}
				}
			}
		}
	}
}
