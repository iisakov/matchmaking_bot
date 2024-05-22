package model

import (
	"errors"
	"strconv"

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
	Text      string
	Markup    tgbotapi.InlineKeyboardMarkup
	MustMatch bool
}

type QuestionsList []Question

var RowQuestion = map[string]map[string][]string{
	"Выбор пола": {
		"answers": {"Вариант ответа 1.1", "Вариант ответа 1.2"},
		"options": {}},
	"Вопрос 2": {
		"answers": {"Вариант ответа 2.1", "Вариант ответа 2.2", "Вариант ответа 2.3", "Вариант ответа 2.4", "Вариант ответа 2.5", "Вариант ответа 2.6"},
		"options": {"mustMatch"}},
	"Вопрос 3": {
		"answers": {"Вариант ответа 3.1", "Вариант ответа 3.2"},
		"options": {}},
	"Вопрос 4": {
		"answers": {"Вариант ответа 4.1", "Вариант ответа 4.2", "Вариант ответа 4.3"},
		"options": {"mustMatch"}},
	"Вопрос 5": {
		"answers": {"Вариант ответа 5.1", "Вариант ответа 5.2", "Вариант ответа 5.3", "Вариант ответа 5.4", "Вариант ответа 5.5", "Вариант ответа 5.6", "Вариант ответа 5.7"},
		"options": {}},
	"Вопрос 6": {
		"answers": {"Вариант ответа 6.1", "Вариант ответа 6.2", "Вариант ответа 6.3", "Вариант ответа 6.4", "Вариант ответа 6.5"},
		"options": {"mustMatch"}},
	"Вопрос 7": {
		"answers": {"Вариант ответа 7.1", "Вариант ответа 7.2", "Вариант ответа 7.3"},
		"options": {}},
	"Вопрос 8": {
		"answers": {"Вариант ответа 8.1", "Вариант ответа 8.2", "Вариант ответа 8.3", "Вариант ответа 8.4"},
		"options": {"mustMatch"}},
	"Вопрос 9": {
		"answers": {"Вариант ответа 9.1", "Вариант ответа 9.2", "Вариант ответа 9.3", "Вариант ответа 9.4"},
		"options": {"mustMatch"}},
	"Вопрос 10": {
		"answers": {"Вариант ответа 10.1", "Вариант ответа 10.2"},
		"options": {"mustMatch"}},
}

func CreateQuestions(rowQuestion map[string]map[string][]string) (result QuestionsList) {

	for questionText, answerValue := range rowQuestion {
		buttons := []tgbotapi.InlineKeyboardButton{}
		for i, a := range answerValue["answers"] {
			buttons = append(buttons, tgbotapi.NewInlineKeyboardButtonData(a, strconv.Itoa(i)))
		}
		keyboardMarkup := tgbotapi.NewInlineKeyboardMarkup(tgbotapi.NewInlineKeyboardRow(buttons...))
		result = append(result, Question{Text: questionText, Markup: keyboardMarkup, MustMatch: findOptions(answerValue["options"], "mustMatch")})
	}
	return
}

var BotQuestions Questions = Questions{QuestionsCounter: 0, QuestionsList: CreateQuestions(RowQuestion)}

func findOptions(ol []string, o string) bool {
	for _, v := range ol {
		if v == o {
			return true
		}
	}
	return false
}
