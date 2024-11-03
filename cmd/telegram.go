package main

import (
	"github.com/digkill/telegram-chatgpt/internal/services/telegram"
	"os"
)

func newTelegram() *telegram.Telegram {
	var token = os.Getenv("TELEGRAM_TOKEN")
	var isDebug = os.Getenv("TELEGRAM_DEBUG") == "true"

	wrapper, err := telegram.NewTelegram(isDebug)
	if err != nil {
		panic(err)
	}

	telegramInstance, err := wrapper.GetInstance(token)
	if err != nil {
		panic(err)
	}

	return telegramInstance

}
