package config

import (
	"errors"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	BotToken     string
	DBURL        string
	YtDlpPath    string
	VideoMaxSize int64
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
	ytdlpPath := os.Getenv("YTDLP_PATH")
	if ytdlpPath == "" {
		ytdlpPath = "./bin/yt-dlp"
	}
	maxSizeStr := os.Getenv("VIDEO_MAX_SIZE")
	var videoMaxSize int64 = 52428800
	if maxSizeStr != "" {
		parsed, err := strconv.ParseInt(maxSizeStr, 10, 64)
		if err == nil {
			videoMaxSize = parsed
		}
	}

	cfg := &Config{
		BotToken:     botToken,
		DBURL:        dbURL,
		YtDlpPath:    ytdlpPath,
		VideoMaxSize: videoMaxSize,
	}

	return cfg, nil
}
