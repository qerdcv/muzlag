package config

import (
	"encoding/json"
	"fmt"
	"os"
)

type Config struct {
	BotToken string `json:"bot_token"`
}

func New(path string) (*Config, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("os open: %w", err)
	}
	defer f.Close()

	var config Config
	if err := json.NewDecoder(f).Decode(&config); err != nil {
		return nil, fmt.Errorf("json decode: %w", err)
	}

	return &config, nil
}
