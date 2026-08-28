package start

import (
	"context"
	"log/slog"

	"BotinokTG/internal/handlers/mainmenu"
	"BotinokTG/internal/users"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

type Handler struct {
	service *users.RegistrationService
}

func NewHandler(service *users.RegistrationService) *Handler {
	return &Handler{
		service: service,
	}
}

func (h *Handler) Start(ctx context.Context, b *bot.Bot, update *models.Update) {
	if update.Message == nil {
		return
	}

	var chatTitle string
	if update.Message.Chat.Type == "private" {
		chatTitle = update.Message.Chat.FirstName
	} else {
		chatTitle = update.Message.Chat.Title
	}

	err := h.service.RegisterUserInChat(
		ctx,
		update.Message.From.ID,
		update.Message.From.Username,
		update.Message.From.FirstName,
		update.Message.Chat.ID,
		chatTitle,
	)
	if err != nil {
		slog.Error("failed to register user in chat", slog.Any("error", err))
	} else {
		slog.Info("user successfully registered in chat")
	}

	// Call InfoHandler to show the welcome message
	mainmenu.InfoHandler(ctx, b, update)
}
