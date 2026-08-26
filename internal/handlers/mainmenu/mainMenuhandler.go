package mainmenu

import (
	"context"
	"log/slog"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

func MainMenuHandler(ctx context.Context, b *bot.Bot, update *models.Update) {
	if update.Message == nil {
		return
	}

	keyboard := &models.InlineKeyboardMarkup{
		InlineKeyboard: [][]models.InlineKeyboardButton{
			{
				{Text: "💰 Finances", CallbackData: "finance"},
				{Text: "🎲 Who today?", CallbackData: "who"},
				{Text: "⏰ Reminder", CallbackData: "reminder"},
				{Text: "📊 Balance", CallbackData: "balance"},
			},
		},
	}

	_, err := b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID:      update.Message.Chat.ID,
		Text:        "Welcome to main menu!",
		ReplyMarkup: keyboard,
	})
	if err != nil {
		slog.Error("failed to send message", slog.Any("error", err))
	}
}
