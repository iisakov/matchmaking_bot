package model

import (
	"errors"

	tgbotapi "github.com/iisakov/telegram-bot-api"
)

type Questions struct {
	QuestionsCounter int
	QuestionsList    QuestionsList
}

func (q Questions) GetCurentQuestion() Question {
	return q.QuestionsList[q.QuestionsCounter]
}

func (q Questions) GetQuestionsCounter() int {
	return q.QuestionsCounter
}

func (q *Questions) Back() (Questions, error) {
	if q.QuestionsCounter-1 > 0 {
		q.QuestionsCounter -= 1
	} else {
		return *q, errors.New("вышли за пределы массива")
	}
	return *q, nil
}

func (q *Questions) Next() (Questions, error) {
	if q.QuestionsCounter+1 < len(q.QuestionsList) {
		q.QuestionsCounter += 1
	} else {
		return *q, errors.New("вышли за пределы массива")
	}
	return *q, nil
}

type Question struct {
	Text     string
	Markup   tgbotapi.InlineKeyboardMarkup
	Keyboard tgbotapi.ReplyKeyboardMarkup
	Options  []string
}

func (q Question) GetAnswers() (result []string) {
	for _, ikbs := range q.Markup.InlineKeyboard {
		for _, ikb := range ikbs {
			result = append(result, *ikb.CallbackData)
		}
	}
	return result
}

type QuestionsList []Question

func (qs QuestionsList) FindQuestionByName(questionName string) *Question {
	for _, q := range qs {
		if q.Text == questionName {
			return &q
		}
	}
	return nil
}

func (qs QuestionsList) IsExistQuestionOptionsByName(questionName, option string) bool {
	for _, q := range qs {
		if q.Text == questionName {
			for _, o := range q.Options {
				if o == option {
					return true
				}
			}
		}
	}
	return false
}

func (qs QuestionsList) GetAnswersByQuestionName(questionName string) []string {
	return qs.FindQuestionByName(questionName).GetAnswers()
}

var RowQuestion = map[string]map[string][]string{
	"Кто вы?": {
		"answers": {"Юноша", "Девушка"},
		"options": {"gender"}},
	"Что вы любите читать? Можно выбрать несколько ответов.": {
		"answers": {"романы", "детективы", "лента инстаграмма", "этикетки"},
		"options": {""}},
	"Какой отдых для вас самый лучший?": {
		"answers": {"на море", "горы", "деревня", "дома полежать/не выходить из дома"},
		"options": {"onlyOne"}},
	"Сколько времени в день вы уделяете работе/учебе?": {
		"answers": {"8 часов", "10 часов", "2-3 часа", "не могу ответить"},
		"options": {"onlyOne"}},
	"Через какое время вы ответите на сообщение вашего партнера?": {
		"answers": {"сразу же", "когда закончу свои дела", "перезвоню", "если важно отвечу сразу"},
		"options": {"onlyOne"}},
	"Кто из родителей занимался домашними обязанностями": {
		"answers": {"мужчина", "женщина", "оба", "наняли помощницу"},
		"options": {"onlyOne"}},
	"Как вы чаще всего/как вы предпочитаете выражать свою любовь ?": {
		"answers": {"дарить подарки", "говорить комплименты", "проводить вместе время", "помогать"},
		"options": {"onlyOne"}},
	"Через какое время вы признаетесь в любви?": {
		"answers": {"Все сразу видно и понятно.", "Не сразу и не буду затягивать.", "Мне нужно больше время", "Признаюсь в ответ."},
		"options": {"onlyOne"}},
	"В какие кружки Вы посещали в детстве?": {
		"answers": {"Танцы", "Музыкальная школа", "Художественная школа", "Спорт"},
		"options": {"onlyOne"}},
	"Если бы завтра был последний день на земле, что бы ты делал в свои последние 24 часа?": {
		"answers": {"Провел день, как обычно", "В кругу семьи", "Тусил", "С любимым человеком"},
		"options": {"onlyOne"}},
}

func CreateQuestionsWithInlineKeyboard(rowQuestion map[string]map[string][]string) (result QuestionsList) {
	for questionText, answerValue := range rowQuestion {
		var rows [][]tgbotapi.InlineKeyboardButton

		for _, a := range answerValue["answers"] {
			rows = append(rows, tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData(a, a)))
		}
		keyboardMarkup := tgbotapi.NewInlineKeyboardMarkup(rows...)
		result = append(result, Question{Text: questionText, Markup: keyboardMarkup, Options: answerValue["options"]})
	}
	return
}

func CreateQuestionsWithKeyboard(rowQuestion map[string]map[string][]string) (result QuestionsList) {
	for questionText, answerValue := range rowQuestion {
		var rows [][]tgbotapi.KeyboardButton

		for _, a := range answerValue["answers"] {
			rows = append(rows, tgbotapi.NewKeyboardButtonRow(tgbotapi.NewKeyboardButton(a)))
		}
		keyboardMarkup := tgbotapi.NewReplyKeyboard(rows...)
		result = append(result, Question{Text: questionText, Keyboard: keyboardMarkup, Options: answerValue["options"]})
	}
	return
}

var BotQuestions Questions = Questions{QuestionsCounter: 0, QuestionsList: CreateQuestionsWithInlineKeyboard(RowQuestion)}
var BotQuestionsK Questions = Questions{QuestionsCounter: 0, QuestionsList: CreateQuestionsWithKeyboard(RowQuestion)}
