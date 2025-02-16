package commands

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"gitlab.com/mediarise/appleclassbot/internal/components/database"
	"gitlab.com/mediarise/appleclassbot/internal/config"
	"gitlab.com/mediarise/appleclassbot/internal/services/telegram"
)

type RefCommandInterface interface {
	Execute(telegramUserId int, botUsername string)
}

type RefCommand struct {
	bot    telegram.Telegram
	config *config.Config
	db     *database.DbComponent
}

func (c *RefCommand) Execute(telegramUserId int, botUsername string) {

	refLink := fmt.Sprintf("https://t.me/%s?start=%d", botUsername, telegramUserId)
	msg := tgbotapi.NewMessage(int64(telegramUserId), "Ваша реферальная ссылка: "+refLink)
	c.bot.Send(msg)
}

func NewRefCommand(bot telegram.Telegram, config *config.Config, db *database.DbComponent) *RefCommand {
	return &RefCommand{
		bot:    bot,
		config: config,
		db:     db,
	}
}
