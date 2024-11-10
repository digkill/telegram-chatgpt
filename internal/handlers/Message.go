package handlers

import (
	"github.com/digkill/telegram-chatgpt/internal/models"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/sirupsen/logrus"
)

type MessageContext struct {
	Updater *UpdateTelegramData
	Payload string
}

type MessageHandler interface {
	Handle(message *tgbotapi.Message, ctx *MessageContext)
}

type InitHandler struct {
	Next MessageHandler
}

func (i *InitHandler) Handle(message *tgbotapi.Message, ctx *MessageContext) {
	if message.IsCommand() && message.Command() == "start" {
		err := ctx.Updater.SendMessageTelegram(
			message.Chat.ID,
			"Привет Чувак!")

		if err != nil {
			logrus.Errorf("Cannot send message. Error: " + err.Error())
		}
	}
	i.Next.Handle(message, ctx)
}

type CommandMenuHandler struct {
	Next MessageHandler
}

func (h *CommandMenuHandler) Handle(message *tgbotapi.Message, ctx *MessageContext) {
	if message.Command() == "start" {
		ctx.Updater.Handler.SendListMenu(
			message.Chat.ID,
			"Привет! Странник!",
			models.Button{
				Type: "show_main_menu",
			},
		)
	} else {
		h.Next.Handle(message, ctx)
	}
}

type FinishHandler struct{}

func (h *FinishHandler) Handle(message *tgbotapi.Message, ctx *MessageContext) {
	err := ctx.Updater.SendMessageTelegram(message.Chat.ID,
		"Проверьте корректность команды!")
	if err != nil {
		logrus.Errorf("Cannot send message. Error: " + err.Error())
	}
}
