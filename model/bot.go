package model

import (
	"errors"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type TgBot struct {
	Bot   *tgbotapi.BotAPI
	Stage BotStage
}

func NewTgBot(bot *tgbotapi.BotAPI) TgBot {
	return TgBot{Bot: bot, Stage: BotStage{StageName: stages[0], StageType: 0, StagesText: stageText[0]}}
}

func (tgb TgBot) SendMsgById(chatId int64, msgText ...string) {
	text := ""
	for _, t := range msgText {
		text = text + "\n" + t
	}
	msg := tgbotapi.NewMessage(chatId, text)
	tgb.Bot.Send(msg)
}

func (tgb TgBot) SendMsgWithInleneKeyboardById(chatId int64, inleneKeyboard tgbotapi.InlineKeyboardMarkup, msgText ...string) {
	text := ""
	for _, t := range msgText {
		text = text + "\n" + t
	}
	msg := tgbotapi.NewMessage(chatId, text)
	msg.ReplyMarkup = inleneKeyboard
	tgb.Bot.Send(msg)
}

type BotStage struct {
	StageName  string
	StageType  int
	StagesText string
}

var stages = []string{"Настройка системы", "Регистрация", "Вопросы", "Распределение пар", "Общение в парах", "Заключение"}
var stageText = []string{
	"В этом режиме, мы настраиваем бота. Проверяем его работоспособность, назначаем модераторов.",
	"Режим регистрации, Пользователи регистрируются на площадке. Придумывают себе псевдонимы.",
	"Режим ответов на вопросы. Только в этом режиме пользователь может ответить на вопросы.",
	"Этот режим нужен для модераторов. Проверить как бот распределил пары, всем ли досталась пара. Нет ли проблем.",
	"Режим приватного общения в парах. Все сообщения перенаправляются через бота собеседникам под псевдонимами.",
	"Прощание и открывание личности собеседника.",
}

func (bs *BotStage) Up() (BotStage, error) {
	if bs.StageType+1 < len(stages) {
		bs.StageType += 1
		bs.StageName = stages[bs.StageType]
		bs.StagesText = stageText[bs.StageType]
	} else {
		return *bs, errors.New("вышли за пределы массива")
	}
	return *bs, nil
}

func (bs *BotStage) Down() (result BotStage, err error) {
	if bs.StageType-1 >= 0 {
		bs.StageType -= 1
		bs.StageName = stages[bs.StageType]
		bs.StagesText = stageText[bs.StageType]
	} else {
		return *bs, errors.New("вышли за пределы массива")
	}
	return *bs, nil
}
