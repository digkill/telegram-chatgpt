package handlers

import (
	"fmt"
	"github.com/digkill/telegram-chatgpt/internal/models"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/sirupsen/logrus"
	"regexp"
)

type CallBackHandler interface {
	Handle(callbackQuery *tgbotapi.CallbackQuery, ctx *CallBackContext)
}

type CallBackContext struct {
	Updater     *UpdateTelegramData
	RequestData *models.Button
	Payload     string
}

type MainMenuHandler struct {
	Next CallBackHandler
}

func (i *MainMenuHandler) Handle(callbackQuery *tgbotapi.CallbackQuery, ctx *CallBackContext) {

	fmt.Println("🪭🪭🪭🪭🪭")
	fmt.Println(ctx.Updater)
	fmt.Println("🪭🪭🪭🪭🪭")

	if ctx.RequestData != nil && ctx.RequestData.Type == "chatGPT" {

		ctx.Updater.Handler.SendListMenu(
			callbackQuery.Message.Chat.ID,
			"Привет! Странник!!",
			models.Button{
				Type: "show_main_menu",
			},
		)
		return
	}
	i.Next.Handle(callbackQuery, ctx)
}

type ChatGPTHandler struct {
	Next CallBackHandler
}

func (i *ChatGPTHandler) Handle(callbackQuery *tgbotapi.CallbackQuery, ctx *CallBackContext) {

	matched, err := regexp.MatchString(`"chatGPT"`, callbackQuery.Data)
	if err != nil {
		logrus.Error(err)
	}

	if matched {

		ctx.Updater.Handler.SendMessageTelegram(
			callbackQuery.Message.Chat.ID,
			"Введите запрос:",
		)
		return

	}

	i.Next.Handle(callbackQuery, ctx)
}

type FinishCallBackHandler struct{}

func (i *FinishCallBackHandler) Handle(callbackQuery *tgbotapi.CallbackQuery, ctx *CallBackContext) {
	err := ctx.Updater.SendMessageTelegram(callbackQuery.Message.Chat.ID,
		"Проверьте корректность команды!")
	if err != nil {
		logrus.Errorf("Cannot send message. Error: " + err.Error())
	}
}
