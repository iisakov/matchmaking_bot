package model

import (
	"strconv"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type Question struct {
	Text   string
	Markup tgbotapi.InlineKeyboardMarkup
}

type Questions []Question

var RowQuestion = map[string][]string{
	"Вопрос 1":  {"Вариант ответа 1.1", "Вариант ответа 1.2", "Вариант ответа 1.3", "Вариант ответа 1.4"},
	"Вопрос 2":  {"Вариант ответа 2.1", "Вариант ответа 2.2", "Вариант ответа 2.3", "Вариант ответа 2.4", "Вариант ответа 2.5", "Вариант ответа 2.6"},
	"Вопрос 3":  {"Вариант ответа 3.1", "Вариант ответа 3.2"},
	"Вопрос 4":  {"Вариант ответа 4.1", "Вариант ответа 4.2", "Вариант ответа 4.3"},
	"Вопрос 5":  {"Вариант ответа 5.1", "Вариант ответа 5.2", "Вариант ответа 5.3", "Вариант ответа 5.4", "Вариант ответа 5.5", "Вариант ответа 5.6", "Вариант ответа 5.7"},
	"Вопрос 6":  {"Вариант ответа 6.1", "Вариант ответа 6.2", "Вариант ответа 6.3", "Вариант ответа 6.4", "Вариант ответа 6.5"},
	"Вопрос 7":  {"Вариант ответа 7.1", "Вариант ответа 7.2", "Вариант ответа 7.3"},
	"Вопрос 8":  {"Вариант ответа 8.1", "Вариант ответа 8.2", "Вариант ответа 8.3", "Вариант ответа 8.4"},
	"Вопрос 9":  {"Вариант ответа 9.1", "Вариант ответа 9.2", "Вариант ответа 9.3", "Вариант ответа 9.4"},
	"Вопрос 10": {"Вариант ответа 10.1", "Вариант ответа 10.2"},
}

func CreateQuestions(rowQuestion map[string][]string) (result Questions) {

	for questionText, answerTexts := range rowQuestion {
		buttons := []tgbotapi.InlineKeyboardButton{}
		for i, a := range answerTexts {
			buttons = append(buttons, tgbotapi.NewInlineKeyboardButtonData(a, strconv.Itoa(i)))
		}
		keyboardMarkup := tgbotapi.NewInlineKeyboardMarkup(tgbotapi.NewInlineKeyboardRow(buttons...))
		result = append(result, Question{Text: questionText, Markup: keyboardMarkup})
	}
	return
}

var QuestionsList = CreateQuestions(RowQuestion)
