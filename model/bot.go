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
	return TgBot{Bot: bot, Stage: BotStage{StageName: "Настройка системы", StageType: 0}}
}

func (tgb TgBot) SendMsgById(chat_id int64, msgText ...string) {
	text := ""
	for _, t := range msgText {
		text = text + "\n" + t
	}
	msg := tgbotapi.NewMessage(chat_id, text)
	tgb.Bot.Send(msg)
}

type BotStage struct {
	StageName string
	StageType int
}

var stages = []string{"Настройка системы", "Регистрация", "Вопросы", "Распределение пар", "Общение в парах", "Заключение"}

func (bs *BotStage) Up() (BotStage, error) {
	if bs.StageType+1 < len(stages) {
		bs.StageType += 1
		bs.StageName = stages[bs.StageType]
	} else {
		return *bs, errors.New("вышли за пределы массива")
	}
	return *bs, nil
}

func (bs *BotStage) Down() (result BotStage, err error) {
	if bs.StageType-1 > 0 {
		bs.StageType -= 1
		bs.StageName = stages[bs.StageType]
	} else {
		return *bs, errors.New("вышли за пределы массива")
	}
	return *bs, nil
}
