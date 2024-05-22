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
	"Выбор пола": {
		"answers": {"Юноша", "Девушка"},
		"options": {"gender"}},
	"Вопрос 2": {
		"answers": {"Вариант ответа 2.1", "Вариант ответа 2.2", "Вариант ответа 2.3", "Вариант ответа 2.4", "Вариант ответа 2.5", "Вариант ответа 2.6"},
		"options": {"onlyOne"}},
	"Вопрос 3": {
		"answers": {"Вариант ответа 3.1", "Вариант ответа 3.2"},
		"options": {"onlyOne"}},
	"Вопрос 4": {
		"answers": {"Вариант ответа 4.1", "Вариант ответа 4.2", "Вариант ответа 4.3"},
		"options": {"onlyOne"}},
	"Вопрос 5": {
		"answers": {"Вариант ответа 5.1", "Вариант ответа 5.2", "Вариант ответа 5.3", "Вариант ответа 5.4", "Вариант ответа 5.5", "Вариант ответа 5.6", "Вариант ответа 5.7"},
		"options": {"onlyOne"}},
	"Вопрос 6": {
		"answers": {"Вариант ответа 6.1", "Вариант ответа 6.2", "Вариант ответа 6.3", "Вариант ответа 6.4", "Вариант ответа 6.5"},
		"options": {"onlyOne"}},
	"Вопрос 7": {
		"answers": {"Вариант ответа 7.1", "Вариант ответа 7.2", "Вариант ответа 7.3"},
		"options": {}},
	"Вопрос 8": {
		"answers": {"Вариант ответа 8.1", "Вариант ответа 8.2", "Вариант ответа 8.3", "Вариант ответа 8.4"},
		"options": {}},
	"Вопрос 9": {
		"answers": {"Вариант ответа 9.1", "Вариант ответа 9.2", "Вариант ответа 9.3", "Вариант ответа 9.4"},
		"options": {}},
	"Вопрос 10": {
		"answers": {"Вариант ответа 10.1", "Вариант ответа 10.2"},
		"options": {}},
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
