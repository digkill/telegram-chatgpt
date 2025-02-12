package main

import (
	"fmt"
	"github.com/digkill/telegram-chatgpt/internal/config"
	"github.com/digkill/telegram-chatgpt/internal/domains"
	"github.com/digkill/telegram-chatgpt/internal/handlers"
	"github.com/digkill/telegram-chatgpt/internal/services/telegram"
	"github.com/sirupsen/logrus"
	"os"
)

func runTelegram() *telegram.Telegram {
	newConfig := config.NewConfig()

	fmt.Println(newConfig)

	wrapper, err := telegram.NewTelegram(newConfig.Telegram.TelegramDebug)
	if err != nil {
		logrus.Panicf("NewTelegram: " + err.Error())
		os.Exit(1)
	}

	bot, err := wrapper.GetInstance(newConfig.Telegram.TelegramToken)
	if err != nil {
		logrus.Panicf("GetInstance: " + err.Error())
		os.Exit(1)
	}

	handler := domains.NewHandler(*bot)

	//	go func() {
	err = handlers.NewUpdateTelegramData(handler).Init()
	if err != nil {
		logrus.Panicf("There is an error in updater telegram data: " + err.Error())
		os.Exit(1)
	}
	//	}()

	return bot

}
