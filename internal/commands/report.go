package commands

import (
	log "github.com/sirupsen/logrus"
	"gitlab.com/mediarise/appleclassbot/internal/components/database"
	"gitlab.com/mediarise/appleclassbot/internal/config"
	"gitlab.com/mediarise/appleclassbot/internal/services/telegram"
)

type ReportCommandInterface interface {
	Execute(telegramUserId int, username string) error
}

type ReportCommand struct {
	bot    telegram.Telegram
	config *config.Config
	db     *database.DbComponent
}

type UserStats struct {
	Date         string
	UserID       int64
	RequestCount int
	TotalCount   int
}

func (s *ReportCommand) Execute() ([]UserStats, error) {

	// SQL-запрос статистики
	query := `
		SELECT 
			DATE(created_at) AS date, 
			user_id, 
			COUNT(*) AS request_count, 
			SUM(count) AS total_count
		FROM journal_tg_gpt
		GROUP BY date, user_id
		ORDER BY date DESC, request_count DESC;
	`

	rows, err := s.db.GetSqlDb().Query(query)
	if err != nil {
		log.Fatalf("Ошибка выполнения запроса: %v", err)
	}
	defer rows.Close()

	// Вывод результатов
	var stats []UserStats
	for rows.Next() {
		var stat UserStats
		if err := rows.Scan(&stat.Date, &stat.UserID, &stat.RequestCount, &stat.TotalCount); err != nil {
			log.Fatalf("Ошибка чтения строки: %v", err)
		}
		stats = append(stats, stat)
	}

	if err := rows.Err(); err != nil {
		log.Fatalf("Ошибка обработки строк: %v", err)
	}

	return stats, nil
}

func NewReportCommand(bot telegram.Telegram, config *config.Config, db *database.DbComponent) *ReportCommand {
	return &ReportCommand{
		bot:    bot,
		config: config,
		db:     db,
	}
}
