package config

import "os"

type Config struct {
	BotToken string `json:"bot_token"`
}

func New() *Config {
	return &Config{
		BotToken: os.Getenv("BOT_TOKEN"),
	}
}
