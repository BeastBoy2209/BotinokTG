package start

import (
	"context"
	"log/slog"

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

	// Определяем название чата
	var chatTitle string
	if update.Message.Chat.Type == "private" {
		chatTitle = update.Message.Chat.FirstName // Для личек используем имя
	} else {
		chatTitle = update.Message.Chat.Title // Для групп используем название
	}

	// Выполняем транзакцию через сервис
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
		// Можно отправить пользователю сообщение об ошибке
	} else {
		slog.Info("user successfully registered in chat")
	}

	keyboard := &models.InlineKeyboardMarkup{
		InlineKeyboard: [][]models.InlineKeyboardButton{
			{
				{Text: "Sign Up", CallbackData: "signup"},
				{Text: "Log In", CallbackData: "login"},
			},
			{
				{Text: "Forgot smth...", CallbackData: "forgot"},
			},
		},
	}

	_, err = b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID:      update.Message.Chat.ID,
		Text:        "Welcome to Botinok",
		ReplyMarkup: keyboard,
	})
	if err != nil {
		slog.Error("smth wrong with message delivery", slog.Any("error", err))
	}
}
