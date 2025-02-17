package commands

import (
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

func (s *StatCommand) Execute(telegramUserId int) (int, error) {
	var count int
	err := s.db.GetSqlDb().QueryRow("SELECT COUNT(*) FROM tg_bot_referral WHERE referrer_id = ?", telegramUserId).Scan(&count)
	if err != nil {
		log.Println("Ошибка при получении статистики:", err)
		return 0, err
	}

	return count, nil

}

func NewStatCommand(bot telegram.Telegram, config *config.Config, db *database.DbComponent) *StatCommand {
	return &StatCommand{
		bot:    bot,
		config: config,
		db:     db,
	}
}
