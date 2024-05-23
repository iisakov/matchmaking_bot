package model

import (
	"errors"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
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
	Text    string
	Markup  tgbotapi.InlineKeyboardMarkup
	Options []string
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
	"Выберите свой пол": {
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
}

func CreateQuestions(rowQuestion map[string]map[string][]string) (result QuestionsList) {
	for questionText, answerValue := range rowQuestion {
		buttons := []tgbotapi.InlineKeyboardButton{}
		for _, a := range answerValue["answers"] {
			buttons = append(buttons, tgbotapi.NewInlineKeyboardButtonData(a, a))
		}
		keyboardMarkup := tgbotapi.NewInlineKeyboardMarkup(tgbotapi.NewInlineKeyboardRow(buttons...))
		result = append(result, Question{Text: questionText, Markup: keyboardMarkup, Options: answerValue["options"]})
	}
	return
}

var BotQuestions Questions = Questions{QuestionsCounter: 0, QuestionsList: CreateQuestions(RowQuestion)}
