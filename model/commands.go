package model

import tgbotapi "github.com/iisakov/telegram-bot-api"

type Comand struct {
	Text     string
	Markup   tgbotapi.InlineKeyboardMarkup
	Keyboard tgbotapi.ReplyKeyboardMarkup
	Options  []string
}

type ComandList []Comand

type Comands struct {
	ComandList QuestionsList
}

func CreateComandsWithKeyboard(rowComands map[string][]string) (result ComandList) {
	for comandText, comand := range rowComands {
		var rows [][]tgbotapi.KeyboardButton

		for _, a := range comand {
			rows = append(rows, tgbotapi.NewKeyboardButtonRow(tgbotapi.NewKeyboardButton(a)))
		}
		keyboardMarkup := tgbotapi.NewReplyKeyboard(rows...)
		result = append(result, Comand{Text: comandText, Keyboard: keyboardMarkup})
	}
	return
}

var rc = map[string][]string{
	"Настройка системы": {"/su", "/sd4"},
	"Регистрация":       {"/gus", "/su", "/sd"},
	"Вопросы":           {"/gus", "/su", "/sd", "/qn", "/qb", "/sq", "/sst"},
	"Распределение пар": {"/su", "/sd1"},
	"Общение в парах":   {"/su", "/sd2"},
	"Заключение":        {"/su", "/sd3"},
}
var C ComandList = CreateComandsWithKeyboard(rc)
