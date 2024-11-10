package handlers

import (
	"fmt"
	"github.com/digkill/telegram-chatgpt/internal/domains"
	"github.com/digkill/telegram-chatgpt/internal/models"
	"github.com/digkill/telegram-chatgpt/internal/services/chatgpt"
	"github.com/gin-gonic/gin"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/sashabaranov/go-openai"
	"github.com/sirupsen/logrus"
	"os"
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

	/*fmt.Println("MainMenuHandler.Handle")
	fmt.Println(callbackQuery.Data)
	fmt.Println("MainMenuHandler.Handle")
	//if ctx.RequestData.Type == "show_main_menu" {

	ctx.Updater.Handler.SendListMenu(
		callbackQuery.Message.Chat.ID,
		"Привет! Странник!!",
		models.Button{
			Type: "show_main_menu",
		},
	)
	//	return
	//} else {
	i.Next.Handle(callbackQuery, ctx)
	//	};;*/
	i.Next.Handle(callbackQuery, ctx)
}

type ChatGPTHandler struct {
	Next CallBackHandler
}

func (i *ChatGPTHandler) Handle(callbackQuery *tgbotapi.CallbackQuery, ctx *CallBackContext) {

	matched, err := regexp.MatchString(`"chatGPT"`, callbackQuery.Data)
	if err != nil {
		fmt.Println("🤡🤡🤡🤡🤡🤡")
		logrus.Error(err)
	}

	fmt.Println("🤡🤡🤡🤡🤡🤡")
	fmt.Println(ctx.RequestData)
	fmt.Println(callbackQuery.Data)
	fmt.Println(matched)
	fmt.Println("🤡🤡🤡🤡🤡🤡")
	if matched {

		var openAIToken = os.Getenv("CHATGPT_TOKEN")
		var openAIURL = os.Getenv("CHATGPT_URL")

		config := openai.DefaultConfig(openAIToken)
		if openAIURL != "" {
			config.BaseURL = openAIURL
		}

		openaiClient := openai.NewClientWithConfig(config)
		chat := chatgpt.NewChatGPT(openaiClient)

		actionInfo := domains.ActionInfo{
			Message: &domains.Message{Role: "User", Content: "Привет лапуля!"},
		}

		msg := actionInfo.GetText()
		messages := domains.MakeMessages(msg)

		var contextGpt *gin.Context

		answer, err := chat.Chat(contextGpt, messages)

		if err != nil {
			logrus.Error(err)
		}

		ctx.Updater.Handler.SendResultAndReturnMenu(
			callbackQuery.Message.Chat.ID,
			answer.Content,
			models.Button{
				Type: "show_main_menu",
			},
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
