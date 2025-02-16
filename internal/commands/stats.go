package commands

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	log "github.com/sirupsen/logrus"
	"gitlab.com/mediarise/appleclassbot/internal/components/database"
	"gitlab.com/mediarise/appleclassbot/internal/config"
	"gitlab.com/mediarise/appleclassbot/internal/services/telegram"
)

type StatCommandInterface interface {
	Execute(telegramUserId int, username string) error
}

type StatCommand struct {
	bot    telegram.Telegram
	config *config.Config
	db     *database.DbComponent
}

func (s *StatCommand) Execute(telegramUserId int, username string) error {
	var count int
	err := s.db.GetSqlDb().QueryRow("SELECT COUNT(*) FROM users WHERE referrer_id = ?", telegramUserId).Scan(&count)
	if err != nil {
		log.Println("Ошибка при получении статистики:", err)
		return err
	}

	msg := tgbotapi.NewMessage(int64(telegramUserId), fmt.Sprintf("У вас %d рефералов! 🎉", count))
	s.bot.Send(msg)
	return nil
}

func NewStatCommand(bot telegram.Telegram, config *config.Config, db *database.DbComponent) *StartCommand {
	return &StartCommand{
		bot:    bot,
		config: config,
		db:     db,
	}
}
