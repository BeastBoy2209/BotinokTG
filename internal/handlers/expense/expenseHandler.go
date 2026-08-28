package expense

import (
	"BotinokTG/internal/expenses"
	"context"
	"fmt"
	"log/slog"
	"strconv"
	"strings"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

type Handler struct {
	service *expenses.ExpenseService
}

func NewHandler(service *expenses.ExpenseService) *Handler {
	return &Handler{service: service}
}

func (h *Handler) HandleEventCommand(ctx context.Context, b *bot.Bot, update *models.Update) {
	if update.Message == nil || update.Message.Text == "" {
		return
	}

	text := update.Message.Text
	parts := strings.Fields(text)

	if len(parts) < 2 {
		h.sendHelp(ctx, b, update.Message.Chat.ID)

		return
	}

	subCmd := parts[1]
	args := parts[2:]
	slog.Info("event command received", slog.String("sub", subCmd), slog.Int("args", len(args)))

	switch subCmd {
	case "create":
		h.createEvent(ctx, b, update, args)
	case "list":
		h.listEvents(ctx, b, update)
	case "history":
		h.showHistory(ctx, b, update)
	case "add":
		h.addExpense(ctx, b, update, args)
	case "debts":
		h.showDebts(ctx, b, update, args)
	case "info":
		h.showInfo(ctx, b, update, args)
	case "close":
		h.closeEvent(ctx, b, update, args)
	default:
		h.sendHelp(ctx, b, update.Message.Chat.ID)
	}
}

func (h *Handler) send(ctx context.Context, b *bot.Bot, chatID int64, text string) {
	_, err := b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: chatID,
		Text:   text,
	})
	if err != nil {
		slog.Error("failed to send message", slog.Any("error", err))
	}
}

func (h *Handler) sendHelp(ctx context.Context, b *bot.Bot, chatID int64) {
	msg := "Управление счетами:\n\n" +
		"/event create [Название] @user1 @user2 — создать счет\n" +
		"Пример: /event create Шашлыки @dan @alex\n\n" +
		"/event list — список активных счетов\n\n" +
		"/event history — история закрытых счетов\n\n" +
		"/event info [ID] — детали счета\n\n" +
		"/event add [ID] [Сумма] [Описание] — добавить расход\n" +
		"Пример: /event add 1 5000 мясо\n\n" +
		"/event debts [ID] — кто кому должен\n\n" +
		"/event close [ID] — закрыть счет"
	h.send(ctx, b, chatID, msg)
}

func (h *Handler) createEvent(
	ctx context.Context,
	b *bot.Bot,
	update *models.Update,
	args []string,
) {
	if len(args) == 0 {
		h.send(
			ctx,
			b,
			update.Message.Chat.ID,
			"Укажите название счета.\nПример: /event create Шашлыки @dan @alex",
		)

		return
	}

	var titleParts []string
	var tagged []string
	for _, arg := range args {
		if strings.HasPrefix(arg, "@") {
			tagged = append(tagged, arg[1:])
		} else {
			titleParts = append(titleParts, arg)
		}
	}

	title := strings.Join(titleParts, " ")
	if title == "" {
		h.send(ctx, b, update.Message.Chat.ID, "Название счета не может быть пустым.")

		return
	}

	slog.Info(
		"creating event",
		slog.String("title", title),
		slog.Int("participants", len(tagged)),
	)

	event, participantNames, err := h.service.CreateEvent(
		ctx,
		update.Message.Chat.ID,
		title,
		tagged,
		update.Message.From.ID,
	)
	if err != nil {
		slog.Error("failed to create event", slog.Any("error", err))
		h.send(ctx, b, update.Message.Chat.ID, "Ошибка: "+err.Error())

		return
	}

	partList := strings.Join(participantNames, ", ")
	msg := fmt.Sprintf(
		"Счет «%s» создан! (ID: %d)\nУчастники: %s\n\nДобавить расход:\n/event add %d [сумма] [описание]",
		event.Title,
		event.ID,
		partList,
		event.ID,
	)
	h.send(ctx, b, update.Message.Chat.ID, msg)
}

func (h *Handler) listEvents(ctx context.Context, b *bot.Bot, update *models.Update) {
	summaries, err := h.service.GetActiveEventSummaries(ctx, update.Message.Chat.ID)
	if err != nil {
		slog.Error("failed to list events", slog.Any("error", err))
		h.send(ctx, b, update.Message.Chat.ID, "Ошибка: "+err.Error())

		return
	}

	if len(summaries) == 0 {
		h.send(
			ctx,
			b,
			update.Message.Chat.ID,
			"Нет активных счетов.\nСоздайте: /event create [Название]",
		)

		return
	}

	msg := "Активные счета:\n\n"
	var msgSb144 strings.Builder
	for _, s := range summaries {
		msgSb144.WriteString(fmt.Sprintf("ID %d — «%s»\n  Участников: %d | Сумма: %d тг\n",
			s.Event.ID, s.Event.Title, s.ParticipantCount, s.TotalAmount/100))
	}
	msg += msgSb144.String()
	msg += "\nПодробности: /event info [ID]"
	h.send(ctx, b, update.Message.Chat.ID, msg)
}

func (h *Handler) showHistory(ctx context.Context, b *bot.Bot, update *models.Update) {
	summaries, err := h.service.GetClosedEventSummaries(ctx, update.Message.Chat.ID)
	if err != nil {
		slog.Error("failed to get history", slog.Any("error", err))
		h.send(ctx, b, update.Message.Chat.ID, "Ошибка: "+err.Error())

		return
	}

	if len(summaries) == 0 {
		h.send(ctx, b, update.Message.Chat.ID, "Закрытых счетов пока нет.")

		return
	}

	msg := "Закрытые счета:\n\n"
	var msgSb166 strings.Builder
	for _, s := range summaries {
		msgSb166.WriteString(fmt.Sprintf("ID %d — «%s»\n  Участников: %d | Сумма: %d тг\n",
			s.Event.ID, s.Event.Title, s.ParticipantCount, s.TotalAmount/100))
	}
	msg += msgSb166.String()
	msg += "\nДетали: /event info [ID] | Долги: /event debts [ID]"
	h.send(ctx, b, update.Message.Chat.ID, msg)
}

func (h *Handler) showInfo(ctx context.Context, b *bot.Bot, update *models.Update, args []string) {
	if len(args) == 0 {
		h.send(ctx, b, update.Message.Chat.ID, "Укажите ID: /event info [ID]")

		return
	}

	eventID, err := strconv.ParseInt(args[0], 10, 64)
	if err != nil {
		h.send(ctx, b, update.Message.Chat.ID, "ID должен быть числом.")

		return
	}

	info, err := h.service.GetEventInfo(ctx, eventID)
	if err != nil {
		slog.Error("failed to get event info", slog.Any("error", err))
		h.send(ctx, b, update.Message.Chat.ID, "Ошибка: "+err.Error())

		return
	}

	msg := fmt.Sprintf("Счет «%s» (ID: %d)", info.Event.Title, info.Event.ID)
	if info.Event.Status == "closed" {
		msg += " — ЗАКРЫТ"
	} else {
		msg += " — открыт"
	}
	msg += "\n\n"
	msg += "Участники: " + strings.Join(info.ParticipantNames, ", ") + "\n\n"

	if len(info.Expenses) == 0 {
		msg += "Расходов пока нет.\n"
	} else {
		msg += "Расходы:\n"
		var msgSb206 strings.Builder
		for _, e := range info.Expenses {
			desc := ""
			if e.Expense.Description.Valid {
				desc = " — " + e.Expense.Description.String
			}
			msgSb206.WriteString(
				fmt.Sprintf(
					"• %s заплатил %d тг%s\n",
					e.PayerName,
					e.Expense.Amount/100,
					desc,
				),
			)
		}
		msg += msgSb206.String()
		msg += fmt.Sprintf("\nИтого: %d тг", info.TotalAmount/100)
	}

	h.send(ctx, b, update.Message.Chat.ID, msg)
}

func (h *Handler) addExpense(
	ctx context.Context,
	b *bot.Bot,
	update *models.Update,
	args []string,
) {
	if len(args) < 2 {
		h.send(ctx, b, update.Message.Chat.ID, "Формат: /event add [ID] [Сумма] [Описание]")

		return
	}

	eventID, err := strconv.ParseInt(args[0], 10, 64)
	if err != nil {
		h.send(ctx, b, update.Message.Chat.ID, "ID должен быть числом.")

		return
	}

	amount, err := strconv.ParseInt(args[1], 10, 64)
	if err != nil || amount <= 0 {
		h.send(ctx, b, update.Message.Chat.ID, "Некорректная сумма.")

		return
	}

	desc := ""
	if len(args) > 2 {
		desc = strings.Join(args[2:], " ")
	}

	eventTitle, err := h.service.AddExpenseToEvent(
		ctx,
		eventID,
		update.Message.From.ID,
		amount*100,
		desc,
	)
	if err != nil {
		slog.Error("failed to add expense", slog.Any("error", err))
		h.send(ctx, b, update.Message.Chat.ID, "Ошибка: "+err.Error())

		return
	}

	descStr := ""
	if desc != "" {
		descStr = " (" + desc + ")"
	}
	msg := fmt.Sprintf(
		"%d тг%s добавлено в счет «%s»\nДолги: /event debts %d",
		amount,
		descStr,
		eventTitle,
		eventID,
	)
	h.send(ctx, b, update.Message.Chat.ID, msg)
}

func (h *Handler) showDebts(ctx context.Context, b *bot.Bot, update *models.Update, args []string) {
	if len(args) == 0 {
		h.send(ctx, b, update.Message.Chat.ID, "Укажите ID: /event debts [ID]")

		return
	}

	eventID, err := strconv.ParseInt(args[0], 10, 64)
	if err != nil {
		return
	}

	result, err := h.service.CalculateDebts(ctx, eventID)
	if err != nil {
		slog.Error("failed to calculate debts", slog.Any("error", err))
		h.send(ctx, b, update.Message.Chat.ID, "Ошибка: "+err.Error())

		return
	}

	msg := fmt.Sprintf(
		"Счет «%s» — %d тг\nУчастников: %d | Доля: %d тг\n\n",
		result.EventTitle,
		result.TotalAmount/100,
		result.ParticipantCount,
		result.SharePerPerson/100,
	)

	if len(result.Debts) == 0 {
		msg += "Все в расчете! Долгов нет."
	} else {
		msg += "Кто кому должен:\n"
		var msgSb282 strings.Builder
		for _, d := range result.Debts {
			msgSb282.WriteString(
				fmt.Sprintf(
					"• %s -> %s: %d тг\n",
					d.DebtorName,
					d.CreditorName,
					d.Amount/100,
				),
			)
		}
		msg += msgSb282.String()
	}

	h.send(ctx, b, update.Message.Chat.ID, msg)
}

func (h *Handler) closeEvent(
	ctx context.Context,
	b *bot.Bot,
	update *models.Update,
	args []string,
) {
	if len(args) == 0 {
		h.send(ctx, b, update.Message.Chat.ID, "Укажите ID: /event close [ID]")

		return
	}

	eventID, err := strconv.ParseInt(args[0], 10, 64)
	if err != nil {
		return
	}

	eventTitle, err := h.service.CloseEvent(ctx, eventID)
	if err != nil {
		slog.Error("failed to close event", slog.Any("error", err))
		h.send(ctx, b, update.Message.Chat.ID, "Ошибка: "+err.Error())

		return
	}

	h.send(ctx, b, update.Message.Chat.ID, fmt.Sprintf("Счет «%s» закрыт!", eventTitle))
}
