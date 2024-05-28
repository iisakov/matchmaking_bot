package tgstl

import (
	"fmt"
	"matchmaking_bot/config"
	"matchmaking_bot/model"

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
		b.SendMsgByIdAndDeleteOtherMsg(
			um.From.ID,
			um.MessageID,
			fmt.Sprintf(
				`Отлично, как бы ты себя не назвал, %s, таким тебя будет видеть собеседник.
У тебя есть ещё пара минут подумать и изменить псевдоним,
просто отправь мне сообщение.`, um.Text))

	case 4: // Общение в парах
		conversationPartnerId, ok := config.PAIRS.GetConversationPartner(um.From.ID)
		if !ok {
			b.SendMsgById(
				config.MODERATOR_BOT_CHAT,
				fmt.Sprintf("Так вышло, что кому-то не досталась пара: %s", config.CUSTOMERS.FindUserById(um.From.ID).UserLogin))
		}
		b.SendMsgById(
			conversationPartnerId,
			fmt.Sprintf("Сообщение от %s:", config.CUSTOMERS.FindUserById(um.From.ID).UserAlias),
			um.Text)
	default:
		b.SendMsgById(
			um.From.ID,
			"Прошу прощения, пока мне нечего на это ответить.")
	}

	config.CUSTOMERS.FindUserByIdAndSetLastMessageId(um.From.ID, um.MessageID)
}

func HandleChannelPostText(um tgbotapi.Message, b model.TgBot) {
	if um.SenderChat.ID == config.MODERATOR_BOT_CHAT {
		b.SendMsgById(
			config.PUBLIC_BOT_CHAT,
			fmt.Sprintf("Сообщение от команды [by_artisan]: %s", um.Text))
		for _, user := range config.CUSTOMERS {
			b.SendMsgById(
				user.UserChat_id,
				fmt.Sprintf("Сообщение от команды [by_artisan]: %s", um.Text))
			b.SendMsgById(
				config.MODERATOR_BOT_CHAT,
				fmt.Sprintf("Пользователю: %s, он же %s, Отправлено сообщение", user.UserLogin, user.UserAlias))
		}
	}
}
