package tgstl

import (
	"matchmaking_bot/config"
	"matchmaking_bot/model"
	"matchmaking_bot/stl"
	"strconv"

	tgbotapi "github.com/iisakov/telegram-bot-api"
)

func HandleMessageCommands(um tgbotapi.Message, b model.TgBot) {
	switch um.Command() {
	case "start":
		if b.Stage.StageType == 1 {
			b.DeleteMessegeByIds(um.From.ID, stl.CreateSliceInt(um.MessageID-100, um.MessageID))
			b.SendMsgById(um.From.ID, "Привет, настало время придумать себе псевдоним.")
			return
		}
		b.DeleteMessegeByIds(um.From.ID, stl.CreateSliceInt(um.MessageID-100, um.MessageID))
		b.SendMsgById(um.From.ID, "Прошу прощения, пока мне нечего на это ответить.")
	default:
		b.DeleteMessegeByIds(um.From.ID, stl.CreateSliceInt(um.MessageID-100, um.MessageID))
		b.SendMsgById(um.From.ID, "Можешь не стараться, кроме команды /start я ничего не знаю.")
	}

	config.CUSTOMERS.FindUserByIdAndSetLastMessageId(um.From.ID, um.MessageID)
}

func HandleChannelPostCommands(um tgbotapi.Message, b model.TgBot) model.TgBot {
	var err error

	switch um.Command() {
	case "help":
		fallthrough
	case "h":
		b.SendMsgById(int64(config.MODERATOR_BOT_CHAT), "Здеь будет подсказка для модераторов")
		return b

	case "chack_stages_text":
		fallthrough
	case "cst":
		b.SendMsgById(int64(config.MODERATOR_BOT_CHAT), b.Stage.StageName, b.Stage.StagesText)
		return b

	case "send_stages_text":
		fallthrough
	case "sst":
		for _, user := range config.CUSTOMERS {
			b.DeleteMessegeByIds(user.UserChat_id, stl.CreateSliceInt(user.LastMessageId-50, user.LastMessageId+50))
			b.SendMsgById(user.UserChat_id, b.Stage.StageName, b.Stage.StagesText)
		}
		b.SendMsgById(int64(config.MODERATOR_BOT_CHAT), b.Stage.StageName, b.Stage.StagesText)
		b.SendMsgById(int64(config.PUBLIC_BOT_CHAT), b.Stage.StageName, b.Stage.StagesText)
		return b

	case "stage_up":
		fallthrough
	case "su":
		if b.Stage.StageType+1 == 4 && len(config.PAIRS) < 1 {
			b.SendMsgById(
				int64(config.MODERATOR_BOT_CHAT),
				"Нельзя переходить к следующему этапу",
				"Распределённых пар:"+strconv.Itoa(len(config.PAIRS)),
			)
			return b
		}
		b.Stage, err = b.Stage.Up()
		if err != nil {
			b.SendMsgById(int64(config.MODERATOR_BOT_CHAT), err.Error(), b.Stage.StageName, strconv.Itoa(b.Stage.StageType))
			return b
		}
		b.SendMsgById(int64(config.MODERATOR_BOT_CHAT), "Мы перешли на следущий этап.", b.Stage.StageName, strconv.Itoa(b.Stage.StageType))
		return b

	case "stage_down":
		fallthrough
	case "sd":
		b.Stage, err = b.Stage.Down()
		if err != nil {
			b.SendMsgById(int64(config.MODERATOR_BOT_CHAT), err.Error(), b.Stage.StageName, strconv.Itoa(b.Stage.StageType))
			return b
		}
		b.SendMsgById(int64(config.MODERATOR_BOT_CHAT), "Мы вернулись на предыдущий этап.", b.Stage.StageName, strconv.Itoa(b.Stage.StageType))
		return b

	case "question_next":
		fallthrough
	case "qn":

		model.BotQuestions, err = model.BotQuestions.Next()
		if err != nil {
			b.SendMsgById(int64(config.MODERATOR_BOT_CHAT), err.Error(), model.BotQuestions.GetCurentQuestion().Text)
			return b
		}
		b.SendMsgById(int64(config.MODERATOR_BOT_CHAT), "Следующий вопрос звучит так.", model.BotQuestions.GetCurentQuestion().Text)
		return b

	case "question_back":
		fallthrough
	case "qb":
		model.BotQuestions, err = model.BotQuestions.Back()
		if err != nil {
			b.SendMsgById(int64(config.MODERATOR_BOT_CHAT), err.Error(), model.BotQuestions.GetCurentQuestion().Text)
			return b
		}
		b.SendMsgById(int64(config.MODERATOR_BOT_CHAT), "Предыдущий вопрос звучит так.", model.BotQuestions.GetCurentQuestion().Text)
		return b

	case "send_question":
		fallthrough
	case "sq":
		if b.Stage.StageType == 2 {
			for _, user := range config.CUSTOMERS {
				b.DeleteMessegeByIds(user.UserChat_id, stl.CreateSliceInt(user.LastMessageId-50, user.LastMessageId+50))
				b.SendMsgWithInleneKeyboardById(
					user.UserChat_id,
					model.BotQuestions.GetCurentQuestion().Markup,
					model.BotQuestions.GetCurentQuestion().Text)
				b.SendMsgById(int64(config.MODERATOR_BOT_CHAT), "Пользователю: "+user.UserLogin+", он же "+user.UserAlias, "Отправлен вопрос.", model.BotQuestions.GetCurentQuestion().Text)
			}
		} else {
			b.SendMsgById(
				int64(config.MODERATOR_BOT_CHAT),
				"Текущий этап - "+b.Stage.StageName+" - Отправлять вопросы пользователям опрометчиво.",
				strconv.Itoa(model.BotQuestions.GetQuestionsCounter()),
				model.BotQuestions.GetCurentQuestion().Text)
		}
		return b

	case "create_pair":
		fallthrough
	case "cp":
		config.PAIRS = model.CheckCompatibility(config.CUSTOMERS.GetUsersByGender(1), config.CUSTOMERS.GetUsersByGender(0))
		b.SendMsgById(
			int64(config.MODERATOR_BOT_CHAT),
			"Создали пары:",
			config.PAIRS.GetPairs(),
		)
		return b

	case "get_pair":
		fallthrough
	case "gp":
		b.SendMsgById(
			int64(config.MODERATOR_BOT_CHAT),
			config.PAIRS.GetPairs(),
		)

		return b

	case "get_users":
		fallthrough
	case "gus":
		if len(config.CUSTOMERS) > 0 {
			b.SendMsgById(
				int64(config.MODERATOR_BOT_CHAT),
				config.CUSTOMERS.GetUsersByGender(1).GetUsers(),
				config.CUSTOMERS.GetUsersByGender(0).GetUsers(),
			)
		} else {
			b.SendMsgById(
				int64(config.MODERATOR_BOT_CHAT),
				"Нет зарегистрированных пользователей.",
			)
		}

		return b
	}

	return b
}
