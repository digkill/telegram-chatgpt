package handlers

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	config "gitlab.com/mediarise/appleclassbot/internal/config"
	"gitlab.com/mediarise/appleclassbot/internal/domains"
	"gitlab.com/mediarise/appleclassbot/internal/services/telegram"
)

type UpdateTelegramData struct {
	*domains.Handler
}

type Updater interface{}

func (updater *UpdateTelegramData) Init() error {
	upd := tgbotapi.NewUpdate(0)
	upd.Timeout = 60

	newConfig := config.NewConfig()

	wrapper, err := telegram.NewTelegram(newConfig.Telegram.TelegramDebug)

	if err != nil {
		panic(err)
	}

	telegramInstance, err := wrapper.GetInstance(newConfig.Telegram.TelegramToken)
	if err != nil {
		panic(err)
	}

	updates, err := telegramInstance.GetUpdatesChan(upd)
	if err != nil {
		return err
	}

	for update := range updates {

		if update.CallbackQuery != nil {
			(
			//&MainMenuHandler{
			//Next:
			&ChatGPTHandler{
				Next: &FinishCallBackHandler{},
				//	},
			}).Handle(update.CallbackQuery, &CallBackContext{
				Updater: updater,
				Payload: "payload",
				Config:  newConfig,
			})
		}

		if update.Message != nil {
			(&InitHandler{
				Next: &CommandMenuHandler{
					Next: &FinishHandler{},
				},
			}).Handle(update.Message, &MessageContext{
				Updater: updater,
				Payload: "payload",
				Config:  newConfig,
			})
		}

	}

	return nil
}

func NewUpdateTelegramData(handler *domains.Handler) *UpdateTelegramData {
	return &UpdateTelegramData{handler}
}
