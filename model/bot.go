package model

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

var tgHostUrl = "https://api.telegram.org/bot"

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

func (tgb TgBot) DeleteMessegeById(chat_id, message_id int) (resp *http.Response, err error) {
	type DM struct {
		Chat_id    int `json:"chat_id"`
		Message_id int `json:"message_id"`
	}

	data, err := json.Marshal(DM{Chat_id: chat_id, Message_id: message_id})
	if err != nil {
		return
	}

	r := bytes.NewReader(data)
	resp, err = http.Post(tgHostUrl+tgb.Bot.Token+"/deleteMessage", "application/json", r)
	if err != nil {
		return
	}
	return
}

func (tgb TgBot) DeleteMessegeByIds(chat_id int, message_ids []int) (result string, err error) {
	type DM struct {
		Chat_id     int   `json:"chat_id"`
		Message_ids []int `json:"message_ids"`
	}
	fmt.Println(message_ids)
	data, err := json.Marshal(DM{Chat_id: chat_id, Message_ids: message_ids})
	if err != nil {
		return
	}

	r := bytes.NewReader(data)
	resp, err := http.Post(tgHostUrl+tgb.Bot.Token+"/deleteMessages", "application/json", r)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	result = string(bodyBytes)

	return
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
