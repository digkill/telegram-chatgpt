package models

import (
	"database/sql"
	"github.com/digkill/telegram-chatgpt/internal/components/database"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"time"
)

type JournalModel struct {
	Id        int64          `db:"id"`
	UserId    int64          `db:"user_id"`
	Prompt    string         `db:"prompt"`
	Count     int            `db:"count"`
	CreatedAt sql.NullString `db:"created_at"`
	UpdatedAt sql.NullString `db:"updated_at"`
}

type Journal struct {
	*database.DbComponent
}

func (journal *Journal) FindJournalByUserId(userId int64) ([]JournalModel, error) {
	var journalModelList []JournalModel
	err := journal.GetSqlDb().Select(&journalModelList, "SELECT * FROM journal_tg_gpt WHERE user_id = ?", userId)

	if err != nil {
		return nil, err
	}
	return journalModelList, nil
}

func (journal *Journal) FindJournalById(id int64) (*JournalModel, error) {
	var journalModel JournalModel
	err := journal.GetSqlDb().QueryRowx("SELECT * FROM journal_tg_gpt WHERE id = ?", id).StructScan(&journalModel)

	if err != nil {
		return nil, err
	}
	return &journalModel, nil
}

func (journal *Journal) CreateJournal(userId int64, prompt string, count int) (*JournalModel, error) {
	transaction, err := journal.GetSqlDb().Begin()
	if err != nil {
		return nil, err
	}

	var result, _ = transaction.Exec("INSERT INTO journal_tg_gpt (user_id, prompt, count, created_at) VALUES (?, ?, ?, ?)", userId, prompt, count, time.Now())
	if result == nil {
		transaction.Rollback()
		return nil, err
	}

	err = transaction.Commit()
	if err != nil {
		return nil, err
	}
	var lastId, _ = result.LastInsertId()

	return journal.FindJournalById(lastId)
}

/*
	func (journal *Journal) UpdateUser(username string) (*JournalModel, error) {
		transaction, err := journal.GetSqlDb().Begin()
		if err != nil {
			return nil, err
		}
		_, err = transaction.Exec("UPDATE users_tg_gpt SET username = $1", username)
		if err != nil {
			transaction.Rollback()
		}
		err = transaction.Commit()
		if err != nil {
			return nil, err
		}
		return journal.FindJournalByUserId(username)
	}

	func (journal *Journal) DeleteUser(username string) error {
		transaction, err := journal.GetSqlDb().Begin()
		if err != nil {
			return err
		}
		_, err = transaction.Exec("DELETE FROM users_tg_gpt WHERE username = $1", username)
		if err != nil {
			transaction.Rollback()
		}
		err = transaction.Commit()
		if err != nil {
			return err
		}
		return nil
	}
*/
func NewJournal(database *database.DbComponent) *Journal {
	return &Journal{database}
}
