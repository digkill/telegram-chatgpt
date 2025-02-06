package models

import (
	"database/sql"
	"fmt"
	database "github.com/digkill/telegram-chatgpt/internal/components"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"time"
)

type UserModel struct {
	Id        int64          `db:"id"`
	Username  string         `db:"username"`
	CreatedAt sql.NullString `db:"created_at"`
	UpdatedAt sql.NullString `db:"updated_at"`
}

type User struct {
	*database.DbComponent
}

func (user *User) FindUserByUsername(username string) (*UserModel, error) {
	var userModel UserModel
	err := user.GetSqlDb().QueryRowx("SELECT * FROM users_tg_gpt WHERE username = ?", username).StructScan(&userModel)
	fmt.Println("@@@@@@@@@@")
	fmt.Println(err)
	fmt.Println("@@@@@@@@@@")
	if err != nil {
		return nil, err
	}
	return &userModel, nil
}

func (user *User) CreateUser(username string) (*UserModel, error) {
	transaction, err := user.GetSqlDb().Begin()
	if err != nil {
		return nil, err
	}

	_, err = transaction.Exec("INSERT INTO users_tg_gpt (username, created_at) VALUES (?, ?)", username, time.Now())
	if err != nil {
		transaction.Rollback()

		return nil, err
	}

	err = transaction.Commit()
	if err != nil {
		return nil, err
	}
	return user.FindUserByUsername(username)
}

func (user *User) UpdateUser(username string) (*UserModel, error) {
	transaction, err := user.GetSqlDb().Begin()
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
	return user.FindUserByUsername(username)
}

func (user *User) DeleteUser(username string) error {
	transaction, err := user.GetSqlDb().Begin()
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

func NewUser(database *database.DbComponent) *User {
	return &User{database}
}
