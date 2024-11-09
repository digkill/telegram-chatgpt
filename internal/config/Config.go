package config

import "os"

type TelegramConfig struct {
	TelegramToken string `mapstructure:"token"`
	TelegramDebug bool   `mapstructure:"debug"`
}

type Config struct {
	Telegram TelegramConfig `mapstructure:"telegram"`
}

func NewConfig() *Config {
	config := &Config{}

	var token = os.Getenv("TELEGRAM_TOKEN")
	var isDebug = os.Getenv("TELEGRAM_DEBUG") == "true"

	config.Telegram = TelegramConfig{
		TelegramToken: token,
		TelegramDebug: isDebug,
	}

	return config
}
