package videohandler

import (
	"BotinokTG/internal/video"
	"context"
	"log/slog"
	"os"
	"time"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

type Handler struct {
	service *video.VideoService
}

func NewHandler(service *video.VideoService) *Handler {
	return &Handler{service: service}
}

func (h *Handler) HandleMessage(ctx context.Context, b *bot.Bot, update *models.Update) {
	if update.Message == nil || update.Message.Text == "" {
		return
	}

	text := update.Message.Text
	chatID := update.Message.Chat.ID

	_, _, supported := video.DetectPlatform(text)
	if !supported {
		return
	}

	loadingMsg, err := b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: chatID,
		Text:   "Видео скачивается, wait...",
		ReplyParameters: &models.ReplyParameters{
			MessageID: update.Message.ID,
		},
	})
	if err != nil {
		slog.Error("failed to send loading message", slog.Any("error", err))
	}

	go func(messageText string, cID int64) {
		defer func() {
			if loadingMsg != nil {
				b.DeleteMessage(ctx, &bot.DeleteMessageParams{
					ChatID:    cID,
					MessageID: loadingMsg.ID,
				})
			}
		}()

		downloadedVideo, err := h.service.ProcessMessage(ctx, messageText)
		if err != nil {
			slog.Error("download failed", slog.Any("error", err))
			b.SendMessage(ctx, &bot.SendMessageParams{
				ChatID: cID,
				Text:   "Ошибка скачивания видео: " + err.Error(),
			})

			return
		}

		if downloadedVideo == nil {
			return
		}

		file, err := os.Open(downloadedVideo.FilePath)
		if err != nil {
			b.SendMessage(ctx, &bot.SendMessageParams{
				ChatID: cID,
				Text:   "Ошибка при чтении файла",
			})

			return
		}
		defer file.Close()
		defer os.Remove(downloadedVideo.FilePath)

		var sendErr error
		for i := range 3 {
			_, sendErr = b.SendVideo(ctx, &bot.SendVideoParams{
				ChatID: cID,
				Video:  &models.InputFileUpload{Filename: "video.mp4", Data: file},
			})

			if sendErr == nil {
				break
			}

			slog.Warn("sendVideo failed, retrying...", slog.Int("attempt", i+1), slog.Any("error", sendErr))
			time.Sleep(2 * time.Second)
			file.Seek(0, 0)
		}

		if sendErr != nil {
			slog.Error("sendVideo error after retries", slog.Any("error", sendErr))
			_ , _ = b.SendMessage(ctx, &bot.SendMessageParams{
				ChatID: cID,
				Text:   "Видео скачалось, но не удалось отправить в telegram",
			})
		}
	}(text, chatID)
}
