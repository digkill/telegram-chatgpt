package commands

import (
	"database/sql"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	log "github.com/sirupsen/logrus"
	"gitlab.com/mediarise/appleclassbot/internal/components/database"
	"gitlab.com/mediarise/appleclassbot/internal/config"
	"gitlab.com/mediarise/appleclassbot/internal/services/telegram"
	"strconv"
	"strings"
)

type StartCommandInterface interface {
	Execute(msg *tgbotapi.Message, telegramUserId int, username string) error
}

type StartCommand struct {
	bot    telegram.Telegram
	config *config.Config
	db     *database.DbComponent
}

func (s *StartCommand) Execute(msg *tgbotapi.Message, telegramUserId int, username string) error {
	parts := strings.Fields(msg.Text)
	var referrerID int

	// Если есть реферальный код
	if len(parts) > 1 {
		refID, err := strconv.Atoi(parts[1])
		if err == nil && refID != telegramUserId {
			referrerID = refID
		}
	}

	// Проверяем, зарегистрирован ли пользователь
	var existingUserID int
	err := s.db.GetSqlDb().QueryRow("SELECT telegram_user_id FROM tg_bot_referral WHERE telegram_user_id = ?", telegramUserId).Scan(&existingUserID)

	fmt.Println(err)

	if err == sql.ErrNoRows {
		// Если нет, регистрируем
		_, err := s.db.GetSqlDb().Exec("INSERT INTO tg_bot_referral (telegram_user_id, referrer_id, username) VALUES (?, ?, ?)", telegramUserId, referrerID, username)
		if err != nil {
			log.Println("Ошибка при регистрации:", err)
			return err
		}

		reply := "Добро пожаловать! 🎉"
		if referrerID > 0 {
			reply += fmt.Sprintf("\nВы зарегистрированы по реферальной ссылке пользователя %d!", referrerID)
		}
		msg := tgbotapi.NewMessage(msg.Chat.ID, reply)
		s.bot.Send(msg)
	} else {
		// Если уже зарегистрирован
		msg := tgbotapi.NewMessage(msg.Chat.ID, "Вы уже зарегистрированы!")
		s.bot.Send(msg)
	}
	return nil
}

func NewStartCommand(bot telegram.Telegram, config *config.Config, db *database.DbComponent) *StartCommand {
	return &StartCommand{
		bot:    bot,
		config: config,
		db:     db,
	}
}
