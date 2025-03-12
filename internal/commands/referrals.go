package commands

import (
	log "github.com/sirupsen/logrus"
	"gitlab.com/mediarise/appleclassbot/internal/components/database"
	"gitlab.com/mediarise/appleclassbot/internal/config"
	"gitlab.com/mediarise/appleclassbot/internal/services/telegram"
)

type ReferralsCommandInterface interface {
	Execute(telegramUserId int, username string) error
}

type ReferralsCommand struct {
	bot    telegram.Telegram
	config *config.Config
	db     *database.DbComponent
}

type ReferralStats struct {
	ReferrerID    int
	ReferrerCount int
}

func (s *ReferralsCommand) Execute() ([]ReferralStats, error) {

	// SQL-запрос статистики
	query := `SELECT tgr.referrer_id AS referrer_id, COUNT(tgr.id) AS referrer_count FROM tg_bot_referral AS tgr WHERE tgr.referrer_id != 0 GROUP BY tgr.referrer_id;`

	rows, err := s.db.GetSqlDb().Query(query)
	if err != nil {
		log.Fatalf("Ошибка выполнения запроса: %v", err)
	}
	defer rows.Close()

	// Вывод результатов
	var rs []ReferralStats
	for rows.Next() {
		var stat ReferralStats
		if err := rows.Scan(&stat.ReferrerID, &stat.ReferrerCount); err != nil {
			log.Fatalf("Ошибка чтения строки: %v", err)
		}
		rs = append(rs, stat)
	}

	if err := rows.Err(); err != nil {
		log.Fatalf("Ошибка обработки строк: %v", err)
	}

	return rs, nil
}

func NewReferralsCommand(bot telegram.Telegram, config *config.Config, db *database.DbComponent) *ReferralsCommand {
	return &ReferralsCommand{
		bot:    bot,
		config: config,
		db:     db,
	}
}
