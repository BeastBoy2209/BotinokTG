package config

import (
	"errors"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	BotToken string
	DBURL    string
}

func Load() (*Config, error) {
	err := godotenv.Load()
	if err != nil {
		return nil, errors.New("cant load .env")
	}
	botToken := os.Getenv("BOTINOK_TOKEN")
	if botToken == "" {
		return nil, errors.New("telegram bot token is empty in .env")
	}
	dbURL := os.Getenv("DB_URL")
	if dbURL == "" {
		return nil, errors.New("DB_URL is empty in .env")
	}

	cfg := &Config{
		BotToken: botToken,
		DBURL:    dbURL,
	}

	return cfg, nil
}
