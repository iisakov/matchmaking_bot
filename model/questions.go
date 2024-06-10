package model

import (
	"encoding/json"
	"errors"
	"os"

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

func LoadQuestionsOnFile(fo string) (result map[string]map[string][]string) {
	f, err := os.ReadFile(fo)
	if err != nil {
		panic(err)
	}

	err = json.Unmarshal(f, &result)
	if err != nil {
		panic(err)
	}

	return
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

var BotQuestions Questions = Questions{QuestionsCounter: 0, QuestionsList: CreateQuestionsWithInlineKeyboard(LoadQuestionsOnFile("questions.json"))}
var BotQuestionsK Questions = Questions{QuestionsCounter: 0, QuestionsList: CreateQuestionsWithKeyboard(LoadQuestionsOnFile("questions.json"))}
