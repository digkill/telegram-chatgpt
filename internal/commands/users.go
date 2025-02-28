package commands

import (
	log "github.com/sirupsen/logrus"
	"gitlab.com/mediarise/appleclassbot/internal/components/database"
	"gitlab.com/mediarise/appleclassbot/internal/config"
	"gitlab.com/mediarise/appleclassbot/internal/services/telegram"
)

type UsersCommandInterface interface {
	Execute(telegramUserId int) error
}

type UsersCommand struct {
	bot    telegram.Telegram
	config *config.Config
	db     *database.DbComponent
}

type UsersStats struct {
	TotalCount int
}

func (s *UsersCommand) Execute() ([]UsersStats, error) {

	// SQL-запрос статистики
	query := `
		SELECT COUNT(u.id) AS total_count FROM users_tg_gpt AS u`

	rows, err := s.db.GetSqlDb().Query(query)
	if err != nil {
		log.Fatalf("Ошибка выполнения запроса: %v", err)
	}
	defer rows.Close()

	// Вывод результатов
	var stats []UsersStats
	for rows.Next() {
		var stat UsersStats
		if err := rows.Scan(&stat.TotalCount); err != nil {
			log.Fatalf("Ошибка чтения строки: %v", err)
		}
		stats = append(stats, stat)
	}

	if err := rows.Err(); err != nil {
		log.Fatalf("Ошибка обработки строк: %v", err)
	}

	return stats, nil
}

func NewUsersCommand(bot telegram.Telegram, config *config.Config, db *database.DbComponent) *UsersCommand {
	return &UsersCommand{
		bot:    bot,
		config: config,
		db:     db,
	}
}
