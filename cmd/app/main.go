package main

import (
	"BotinokTG/internal/chats"
	"BotinokTG/internal/config"
	"BotinokTG/internal/expenses"
	"BotinokTG/internal/handlers"
	"BotinokTG/internal/members"
	"BotinokTG/internal/storage"
	"BotinokTG/internal/users"
	"BotinokTG/internal/video"
	"BotinokTG/internal/who"
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

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

	// SIGTERM
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	startService := users.NewRegistrationService(dbpool)

	b, err := bot.New(
		cfg.BotToken,
		bot.WithMiddlewares(handlers.RegistrationMiddleware(startService)),
	)
	if err != nil {
		slog.Error("smth went wrong with bot", slog.Any("error", err))
		os.Exit(1)
	}

	userRepo := users.NewRepository(dbpool)
	chatRepo := chats.NewRepository(dbpool)
	memberRepo := members.NewRepository(dbpool)
	expenseRepo := expenses.NewRepository(dbpool)

	expenseService := expenses.NewService(expenseRepo, memberRepo, userRepo, chatRepo)

	ytDlp := video.NewYTDLPDownloader(cfg.YtDlpPath)
	tkDlp := video.NewTikTokDownloader()
	videoService := video.NewService(ytDlp, tkDlp, cfg.VideoMaxSize)

	whoService := who.NewRandomizerService(memberRepo)

	handlers.RegisterAll(b, startService, expenseService, videoService, whoService, memberRepo)

	slog.Info("bot started... \n")

	b.Start(ctx)
	slog.Info("graceful shutdown: ok")
}
