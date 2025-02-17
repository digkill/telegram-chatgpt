package handlers

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/sirupsen/logrus"
	"gitlab.com/mediarise/appleclassbot/internal/commands"
	"gitlab.com/mediarise/appleclassbot/internal/components/database"
	"gitlab.com/mediarise/appleclassbot/internal/config"
	"gitlab.com/mediarise/appleclassbot/internal/models"
)

type CallBackHandler interface {
	Handle(callbackQuery *tgbotapi.CallbackQuery, ctx *CallBackContext)
}

type CallBackContext struct {
	Updater     *UpdateTelegramData
	RequestData *models.Button
	Payload     string
	Config      *config.Config
}

type MainMenuHandler struct {
	Next CallBackHandler
}

func (i *MainMenuHandler) Handle(callbackQuery *tgbotapi.CallbackQuery, ctx *CallBackContext) {

	i.Next.Handle(callbackQuery, ctx)
}

type ChatGPTHandler struct {
	Next CallBackHandler
}

func (i *ChatGPTHandler) Handle(callbackQuery *tgbotapi.CallbackQuery, ctx *CallBackContext) {

	var databaseConfig = ctx.Config.DB
	var db = database.NewDb(&databaseConfig)

	if callbackQuery.Data == "ref" {
		me, err := ctx.Updater.Handler.GetBot().GetMe()
		if err != nil {
			fmt.Println(err)
		}

		refLink := fmt.Sprintf("https://t.me/%s?start=%d", me.UserName, callbackQuery.From.ID)
		msg := tgbotapi.NewMessage(int64(callbackQuery.Message.From.ID), fmt.Sprintf("Ваша реферальная ссылка: [%s](%s)", refLink, refLink))
		err = ctx.Updater.Handler.SendMessageTelegram(
			callbackQuery.Message.Chat.ID,
			msg.Text,
		)
		if err != nil {
			fmt.Println(err)
		}
		return
	}

	if callbackQuery.Data == "stats" {
		statsCommand := commands.NewStatCommand(ctx.Updater.Handler.GetBot(), ctx.Config, db)
		count, err := statsCommand.Execute(callbackQuery.From.ID)
		if err != nil {
			fmt.Println(err)
		}

		msg := tgbotapi.NewMessage(int64(callbackQuery.From.ID), fmt.Sprintf("У вас %d рефералов! 🎉", count))
		err = ctx.Updater.Handler.SendMessageTelegram(
			callbackQuery.Message.Chat.ID,
			msg.Text,
		)
		if err != nil {
			fmt.Println(err)
		}

		return
	}

	i.Next.Handle(callbackQuery, ctx)
}

type RefHandler struct {
	Next CallBackHandler
}

func (i *RefHandler) Handle(callbackQuery *tgbotapi.CallbackQuery, ctx *CallBackContext) {

	if callbackQuery.Data == "ref_menu" {

		me, err := ctx.Updater.GetBot().GetMe()
		if err != nil {
			fmt.Println(err)
		}

		ctx.Updater.Handler.SendRefMenu(
			callbackQuery.Message.Chat.ID,
			fmt.Sprintf("💌 Вы можете пригласить друзей и получить дополнительно 10 запросов в день за каждого друга!\n\n- Когда ваш друг запустит бота, вы получите дополнительно 10 запросов в день;\n- Вы можете пригласить неограниченное количество друзей;\n- Ваш друг должен впервые воспользоваться ботом по вашей персональной ссылке;\n\nСсылка (скопируй ее и отправь другу):  https://t.me/%s?start=%d\n\nИли просто перешлите сообщение ниже своим друзьям:", me.UserName, callbackQuery.From.ID),
			models.Button{
				Type: "show_main_menu",
			},
		)
		err = ctx.Updater.Handler.SendMessageTelegram(
			callbackQuery.Message.Chat.ID,
			fmt.Sprintf("Вы приглашены в бота [%s](https://t.me/%s?start=%d)!\nНажмите на ссылку, чтобы начать:\n🚀 [Запустить бота](https://t.me/%s?start=%d)",
				me.UserName, me.UserName, callbackQuery.From.ID, me.UserName, callbackQuery.From.ID,
			),
		)
		if err != nil {
			fmt.Println(err)
		}
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
