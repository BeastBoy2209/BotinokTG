package expense

import (
	"context"
	"log/slog"
	"strings"
	"strconv"
	"fmt"

	"BotinokTG/internal/expenses"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

type Handler struct{
	service *expenses.ExpenseService
}

func NewHandler(service *expenses.ExpenseService) *Handler{
	return &Handler{service: service}
}

func (h *Handler) AddExpenseHandler(ctx context.Context, b *bot.Bot, update *models.Update) {
	text := update.Message.Text
	parts := strings.SplitN(text, " ", 3)
	if parts[0] != "/expense" {
		return // Игнорируем, если это была другая команда (например /expenses)
	}
	if len(parts) < 2 {
		b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: update.Message.Chat.ID,
			Text:   "укажи сумму и описание",
		})
		return
	}

	amount, err := strconv.ParseInt(parts[1], 10, 64)
	if err != nil || amount <= 0 {
		b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: update.Message.Chat.ID,
			Text:   "число не может быть меньше 0",
		})
		return
	}
	amountInTiyns := amount * 100
	var description string
	if len(parts) == 3 {
		description = parts[2]
	}

	_, err = h.service.AddExpense(ctx, update.Message.Chat.ID, update.Message.From.ID, amountInTiyns, description)
	if err != nil {
		b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: update.Message.Chat.ID,
			Text:   err.Error(),
		})
		return
	}

	slog.Info("expense added", slog.Int64("amount", amount))

	b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: update.Message.Chat.ID,
		Text:   fmt.Sprintf("💰 Расход добавлен\n%d ₸ — %s", amount, description),
	})
}

func (h *Handler) GetExpensesHandler(ctx context.Context, b *bot.Bot, update *models.Update) {
	expensesList, err := h.service.GetChatExpenses(ctx, update.Message.Chat.ID)
	if err != nil {
		b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: update.Message.Chat.ID,
			Text:   err.Error(),
		})
		return
	}

	messageText := "💰 Расходы чата:\n\n"
	for _, exp := range expensesList {
		messageText += fmt.Sprintf("• %d ₸ — %s — %s\n", exp.Amount/100, exp.Description.String, exp.FirstName.String)
	}

	b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: update.Message.Chat.ID,
		Text:   messageText,
		})
}