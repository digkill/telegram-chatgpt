package domains

import (
	"encoding/json"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/sirupsen/logrus"
	"gitlab.com/mediarise/appleclassbot/internal/config"
	"gitlab.com/mediarise/appleclassbot/internal/models"
	"gitlab.com/mediarise/appleclassbot/internal/services/telegram"
)

type Handler struct {
	bot    telegram.Telegram
	Config *config.Config
}

func (handler *Handler) GetBot() telegram.Telegram {
	return handler.bot
}

func (handler *Handler) SendMessageTelegram(chatId int64, message string) error {
	msg := tgbotapi.NewMessage(chatId, message)
	msg.ParseMode = tgbotapi.ModeHTML

	return handler.SendMessageObjectTelegram(msg)
}

func (handler *Handler) SendMessageTextTelegram(chatId int64, message string, parseMode string) (tgbotapi.Message, error) {
	msg := tgbotapi.NewMessage(chatId, message)

	msg.ParseMode = parseMode
	if parseMode == "" {
		msg.ParseMode = tgbotapi.ModeMarkdown
	}
	return handler.bot.Send(msg)

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
			logrus.Errorf("There is an error in send message to user. [SendMessageObjectTelegram] Error: " + err.Error() + " Message.go: " + copyMessageObj.Text)
			return err
		}
	}

	return nil
}

func (handler *Handler) SendMainMenu(chatId int64, message string, data models.Button) bool {
	err := handler.SendMessageWithButtonsInRowToTelegram(
		chatId,
		message,
		tgbotapi.NewInlineKeyboardButtonData("Меню", handler.ButtonToString(data)),
	)
	if err == nil {
		return true
	}
	return false
}

func (handler *Handler) SendListMenu(chatId int64, message string, data models.Button) bool {
	data.Type = "chatGPT"
	chatGPTButton := handler.ButtonToString(data)
	//data.Type = "baton"
	//buttonButton := handler.buttonToString(data)
	err := handler.SendMessageWithButtonsInRowsToTelegram(
		chatId,
		message,
		tgbotapi.NewInlineKeyboardMarkup(
			tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonData("ChatGPT", chatGPTButton),
				//	tgbotapi.NewInlineKeyboardButtonData("Кнопка", buttonButton),
			),
			tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonURL("Подробнее",
					"https://mediarise.org"),
			),
		),
	)
	if err == nil {
		return true
	}
	return false
}
func (handler *Handler) SendRefMenu(chatId int64, message string, data models.Button) bool {

	err := handler.SendMessageWithButtonsInRowsToTelegram(
		chatId,
		message,
		tgbotapi.NewInlineKeyboardMarkup(
			tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonData("⛓️‍💥 Реферальная ссылка", "ref"),
				tgbotapi.NewInlineKeyboardButtonData("👥 Количество приглашенных", "stats"),
			),
		),
	)
	if err == nil {
		return true
	}
	return false
}

func (handler *Handler) SendResultAndReturnMenu(chatId int64, message string, data models.Button) bool {
	err := handler.SendMessageWithButtonsInRowToTelegram(
		chatId,
		message,
		tgbotapi.NewInlineKeyboardButtonData("Вернуться в меню", handler.ButtonToString(data)),
	)
	if err == nil {
		return true
	}
	return false
}

func (handler *Handler) SendResult(chatId int64, message string, data models.Button) bool {
	err := handler.SendMessageTelegram(
		chatId,
		message,
	)
	if err == nil {
		return true
	}
	return false
}

func (handler *Handler) SendMessageWithButtonsInRowToTelegram(chatId int64, message string, buttons ...tgbotapi.InlineKeyboardButton) error {
	return handler.SendMessageWithButtonsInRowsToTelegram(chatId, message, tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(buttons...),
	))
}

func (handler *Handler) SendMessageWithButtonsInRowsToTelegram(chatId int64, message string, markup tgbotapi.InlineKeyboardMarkup) error {
	msg := tgbotapi.NewMessage(chatId, message)
	msg.ParseMode = tgbotapi.ModeHTML
	msg.ReplyMarkup = &markup
	return handler.SendMessageObjectTelegram(msg)
}

func (handler *Handler) ButtonToString(data models.Button) string {
	result, err := json.Marshal(data)
	if err == nil {
		return string(result)
	}
	return ""
}

func (handler *Handler) RemoveMessage(chatId int64, messageId int) (tgbotapi.APIResponse, error) {
	deleteMessageConfig := tgbotapi.NewDeleteMessage(chatId, messageId)

	response, err := handler.bot.DeleteMessage(deleteMessageConfig)
	if err != nil {
		return tgbotapi.APIResponse{}, err
	}
	return response, nil

}

func NewHandler(bot telegram.Telegram) *Handler {
	return &Handler{bot: bot}
}
