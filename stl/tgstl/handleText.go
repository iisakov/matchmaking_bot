package tgstl

import (
	"fmt"
	"matchmaking_bot/config"
	"matchmaking_bot/model"
	"matchmaking_bot/stl"

	tgbotapi "github.com/iisakov/telegram-bot-api"
)

func HandleMessagesText(um tgbotapi.Message, b model.TgBot) {
	switch b.Stage.StageType {
	case 1: // Регистрация
		if !config.CUSTOMERS.IsExistUserById(um.From.ID) {
			config.CUSTOMERS = append(config.CUSTOMERS, model.NewUser(um))
		} else {
			config.CUSTOMERS.FindUserByIdAndUpdateAlias(um.From.ID, um.Text)
		}
		b.DeleteMessegeByIds(um.From.ID, stl.CreateSliceInt(um.MessageID-100, um.MessageID))
		b.SendMsgById(um.From.ID, "Отлично, как бы ты себя не назвал, "+um.Text+", таким тебя будет видеть собеседник. У тебя есть ещё пара минут подумать и изменить псевдоним, просто отправь мне сообщение.")
	case 4: // Общение в парах
		conversationPartnerId, ok := config.PAIRS.GetConversationPartner(um.From.ID)
		if !ok {
			b.SendMsgById(int64(config.MODERATOR_BOT_CHAT), fmt.Sprintf("Так вышло, что кому-то не досталась пара: %s", config.CUSTOMERS.FindUserById(um.From.ID).UserLogin))
		}
		b.SendMsgById(
			conversationPartnerId,
			"Сообщение от "+config.CUSTOMERS.FindUserById(um.From.ID).UserAlias+":",
			um.Text)
	default:
		b.DeleteMessegeByIds(um.From.ID, stl.CreateSliceInt(um.MessageID-100, um.MessageID))
		b.SendMsgById(um.From.ID, "Прошу прощения, пока мне нечего на это ответить.")
	}

	config.CUSTOMERS.FindUserByIdAndSetLastMessageId(um.From.ID, um.MessageID)
}

func HandleChannelPostText(um tgbotapi.Message, b model.TgBot) {
	if um.SenderChat.ID == int64(config.MODERATOR_BOT_CHAT) {
		b.SendMsgById(int64(config.PUBLIC_BOT_CHAT), "Сообщение от команды [by_artisan]:", um.Text)
		for _, user := range config.CUSTOMERS {
			b.SendMsgById(int64(user.UserChat_id), "Сообщение от команды [by_artisan]:", um.Text)
			b.SendMsgById(int64(config.MODERATOR_BOT_CHAT), "Пользователю: "+user.UserLogin+", он же "+user.UserAlias, "Отправлено сообщение")
		}
	}
}
