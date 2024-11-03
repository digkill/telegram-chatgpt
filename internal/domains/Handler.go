package domains

import (
	"github.com/digkill/telegram-chatgpt/internal/services/telegram"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/sirupsen/logrus"
)

type Handler struct {
	bot telegram.Telegram
}

func (handler *Handler) SendMessageTelegram(chatId int64, message string) error {
	msg := tgbotapi.NewMessage(chatId, message)
	msg.ParseMode = tgbotapi.ModeMarkdown

	return handler.SendMessageObjectTelegram(msg)
}

func (handler *Handler) SendMessageObjectTelegram(message tgbotapi.MessageConfig) error {
	runes := []rune(message.Text)

	for i := 0; i < len(runes); i += 4096 {
		nn := i + 4096
		if nn > len(runes) {
			nn = len(runes)
		}

		copyMessageObj := message
		copyMessageObj.Text = string(runes[i:nn])

		_, err := handler.bot.Send(copyMessageObj)
		if err != nil {
			logrus.Errorf("There is an error in send message to user. [SendMessageObjectTelegram] Error: " + err.Error() + " Message: " + copyMessageObj.Text)
			return err
		}
	}

	return nil
}

func NewHandler(bot telegram.Telegram) *Handler {
	return &Handler{bot: bot}
}
