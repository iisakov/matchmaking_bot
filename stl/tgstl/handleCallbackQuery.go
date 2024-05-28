package tgstl

import (
	"matchmaking_bot/config"
	"matchmaking_bot/model"

	tgbotapi "github.com/iisakov/telegram-bot-api"
)

func HandleCallbackQuery(ucq tgbotapi.CallbackQuery, b model.TgBot) {
	if b.Stage.StageType == 2 {
		if config.CUSTOMERS.IsExistUserById(ucq.From.ID) {
			switch {
			case model.BotQuestions.QuestionsList.IsExistQuestionOptionsByName(ucq.Message.Text, "gender"):
				config.CUSTOMERS.FindUserByIdSetGender(ucq.From.ID, ucq.Data)
				b.SendMsgById(
					ucq.From.ID,
					"Отлично, теперь мы знаем какого ты пола.")
				return

			case model.BotQuestions.QuestionsList.IsExistQuestionOptionsByName(ucq.Message.Text, "onlyOne"):
				config.CUSTOMERS.FindUserByIdAndUpdateAnswer(
					ucq.From.ID,
					model.BotQuestions.QuestionsList.GetAnswersByQuestionName(ucq.Message.Text),
					ucq.Data)
				b.SendMsgById(
					ucq.From.ID,
					"Отлично, Можешь изменить ответ если хочешь.")
				return

			default:
				config.CUSTOMERS.FindUserByIdAndAddAnswer(ucq.From.ID, ucq.Data)
				b.SendMsgById(
					ucq.From.ID,
					"Отлично, На этот вопрос можно ответить несколько раз, Выбирай хоть все варианты.")
				return
			}
		}
	}
}
