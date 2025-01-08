package handlers

import (
	"github.com/digkill/telegram-chatgpt/internal/domains"
	"github.com/digkill/telegram-chatgpt/internal/models"
	"github.com/digkill/telegram-chatgpt/internal/services/chatgpt"
	"github.com/gin-gonic/gin"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/sashabaranov/go-openai"
	"github.com/sirupsen/logrus"
	"os"
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
			"Выберите услугу",
			models.Button{
				Type: "show_main_menu",
			},
		)

		return
	} else {

		var openAIToken = os.Getenv("CHATGPT_TOKEN")
		var openAIURL = os.Getenv("CHATGPT_URL")

		config := openai.DefaultConfig(openAIToken)
		if openAIURL != "" {
			config.BaseURL = openAIURL
		}

		openaiClient := openai.NewClientWithConfig(config)
		chat := chatgpt.NewChatGPT(openaiClient)

		actionInfo := domains.ActionInfo{
			Message: &domains.Message{Role: "User", Content: message.Text},
		}

		msg := actionInfo.GetText()
		messages := domains.MakeMessages(msg)

		var contextGpt *gin.Context
		contextGpt = &gin.Context{}

		answer, err := chat.Chat(contextGpt, messages)

		if err != nil {
			logrus.Error(err)
		}

		ctx.Updater.Handler.SendResultAndReturnMenu(
			message.Chat.ID,
			answer.Content,
			models.Button{
				Type: "show_main_menu",
			},
		)
		return
	}

	h.Next.Handle(message, ctx)
}

type FinishHandler struct{}

func (h *FinishHandler) Handle(message *tgbotapi.Message, ctx *MessageContext) {
	err := ctx.Updater.SendMessageTelegram(message.Chat.ID,
		"Проверьте корректность команды!")
	if err != nil {
		logrus.Errorf("Cannot send message. Error: " + err.Error())
	}
}
