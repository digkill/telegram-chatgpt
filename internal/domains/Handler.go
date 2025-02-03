package domains

import (
	"encoding/json"
	"github.com/digkill/telegram-chatgpt/internal/config"
	"github.com/digkill/telegram-chatgpt/internal/models"
	"github.com/digkill/telegram-chatgpt/internal/services/telegram"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/sirupsen/logrus"
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
		tgbotapi.NewInlineKeyboardButtonData("Меню", handler.buttonToString(data)),
	)
	if err == nil {
		return true
	}
	return false
}

func (handler *Handler) SendListMenu(chatId int64, message string, data models.Button) bool {
	data.Type = "chatGPT"
	chatGPTButton := handler.buttonToString(data)
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

func (handler *Handler) SendResultAndReturnMenu(chatId int64, message string, data models.Button) bool {
	err := handler.SendMessageWithButtonsInRowToTelegram(
		chatId,
		message,
		tgbotapi.NewInlineKeyboardButtonData("Вернуться в меню", handler.buttonToString(data)),
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

func (handler *Handler) buttonToString(data models.Button) string {
	result, err := json.Marshal(data)
	if err == nil {
		return string(result)
	}
	return ""
}

func NewHandler(bot telegram.Telegram) *Handler {
	return &Handler{bot: bot}
}
