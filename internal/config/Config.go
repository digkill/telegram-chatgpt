package config

import (
	"os"
	"strconv"
)

type TelegramConfig struct {
	TelegramToken string          `mapstructure:"token"`
	TelegramDebug bool            `mapstructure:"debug"`
	Migration     MigrationConfig `mapstructure:"migration"`
}

type DatabaseConfig struct {
	Type         string `mapstructure:"type"`
	Host         string `mapstructure:"host"`
	Port         int    `mapstructure:"port"`
	Username     string `mapstructure:"username"`
	Password     string `mapstructure:"password"`
	Name         string `mapstructure:"name"`
	MaxIdleConns int    `mapstructure:"max_idle"`
	MaxOpenConns int    `mapstructure:"max_open"`
}

type MigrationConfig struct {
	Path           string `mapstructure:"path"`
	MigrationTable string `mapstructure:"migration_table"`
	DatabaseName   string `mapstructure:"database_name"`
}

type Config struct {
	Telegram TelegramConfig `mapstructure:"telegram"`
	DB       DatabaseConfig `mapstructure:"database"`
}

func NewConfig() *Config {
	config := &Config{}

	var token = os.Getenv("TELEGRAM_TOKEN")
	var isDebug = os.Getenv("TELEGRAM_DEBUG") == "true"

	config.Telegram = TelegramConfig{
		TelegramToken: token,
		TelegramDebug: isDebug,
	}

	var dbHost = os.Getenv("DB_HOST")
	var dbPort, _ = strconv.Atoi(os.Getenv("DB_PORT"))
	var dbUser = os.Getenv("DB_USER")
	var dbPass = os.Getenv("DB_PASS")
	var dbName = os.Getenv("DB_NAME")
	var dbType = os.Getenv("DB_TYPE")
	var maxIdleConns, _ = strconv.Atoi(os.Getenv("MAX_IDLE_CONNS"))
	var maxOpenConns, _ = strconv.Atoi(os.Getenv("MAX_OPEN_CONNS"))

	config.DB = DatabaseConfig{
		Type:         dbType,
		Host:         dbHost,
		Port:         dbPort,
		Username:     dbUser,
		Password:     dbPass,
		Name:         dbName,
		MaxIdleConns: maxIdleConns,
		MaxOpenConns: maxOpenConns,
	}

	return config
}
