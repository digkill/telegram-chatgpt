package commands

import tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"

type CommandInterface interface {
	Execute(msg *tgbotapi.Message, telegramUserId int, username string) error
}
