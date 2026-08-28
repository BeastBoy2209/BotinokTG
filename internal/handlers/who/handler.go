package whohandler

import (
	"BotinokTG/internal/who"
	"context"
	"crypto/rand"
	"math/big"
	"strings"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

type Handler struct {
	service *who.RandomizerService
}

func NewHandler(service *who.RandomizerService) *Handler {
	return &Handler{
		service: service,
	}
}

func (h *Handler) HandleWho(ctx context.Context, b *bot.Bot, update *models.Update) {
	if update.Message == nil || update.Message.Text == "" {
		return
	}

	chatID := update.Message.Chat.ID
	parts := strings.Fields(update.Message.Text)

	// Если переданы аргументы после команды (например, /who @user1 @user2)
	if len(parts) > 1 {
		candidates := parts[1:]

		n, err := rand.Int(rand.Reader, big.NewInt(int64(len(candidates))))
		if err != nil {
			b.SendMessage(ctx, &bot.SendMessageParams{
				ChatID: chatID,
				Text:   "Ошибка рандомизатора.",
			})

			return
		}

		winner := candidates[n.Int64()]
		b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: chatID,
			Text:   "Рандом выбрал: " + winner + "!",
		})

		return
	}

	member, err := h.service.GetRandomMember(ctx, chatID)
	if err != nil {
		if err.Error() == "not enough members" {
			b.SendMessage(ctx, &bot.SendMessageParams{
				ChatID: chatID,
				Text:   "Слишком мало участников для выбора",
			})

			return
		}

		b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: chatID,
			Text:   "Произогла ошибка при выборе",
		})

		return
	}

	var displayName string
	if member.Username.Valid && member.Username.String != "" {
		displayName = "@" + member.Username.String
	} else {
		displayName = member.FirstName
	}

	b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: chatID,
		Text:   "Рандом выбрал " + displayName + "!",
	})
}
