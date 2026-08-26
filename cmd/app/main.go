package main

import (
	"context"
	"log/slog"
	"os"
	
	"BotinokTG/internal/chats"
	"BotinokTG/internal/config"
	"BotinokTG/internal/handlers"
	"BotinokTG/internal/storage"
	"BotinokTG/internal/users"
	"BotinokTG/internal/members"
	"BotinokTG/internal/expenses"
	"BotinokTG/internal/video"

	"github.com/go-telegram/bot"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		slog.Error("failed to load config", slog.Any("error", err))
		os.Exit(1)
	}
	slog.Info("config was loaded... \n")

	dbpool, err := storage.NewPostgres(context.Background(), cfg.DBURL)
	if err != nil {
		slog.Error("can't ping db...")
		os.Exit(1)
	}
	defer dbpool.Close()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	b, err := bot.New(cfg.BotToken)
	if err != nil {
		slog.Error("smth went wrong with bot", slog.Any("error", err))
		os.Exit(1)
	}
	slog.Info("bot started... \n")

	userRepo := users.NewRepository(dbpool)
	chatRepo := chats.NewRepository(dbpool)
	memberRepo := members.NewRepository(dbpool)
	expenseRepo := expenses.NewRepository(dbpool)

	expenseService := expenses.NewService(expenseRepo, memberRepo, userRepo, chatRepo)
	startService := users.NewRegistrationService(dbpool)

	ytDlp := video.NewYTDLPDownloader(cfg.YtDlpPath)
	tkDlp := video.NewTikTokDownloader()
	videoService := video.NewService(ytDlp, tkDlp, cfg.VideoMaxSize)

	handlers.RegisterAll(b, startService, expenseService, videoService)


	b.Start(ctx)
}
