package handlers

import (
	// "context"
	// "log/slog"

	"BotinokTG/internal/expenses"
	"BotinokTG/internal/handlers/expense"
	"BotinokTG/internal/handlers/mainmenu"
	"BotinokTG/internal/handlers/start"
	"BotinokTG/internal/users"
	"BotinokTG/internal/handlers/videohandler"
	"BotinokTG/internal/video"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

func RegisterAll(b *bot.Bot, startService *users.RegistrationService, expenseService *expenses.ExpenseService, videoService *video.VideoService) {
	startHandler := start.NewHandler(startService)

	b.RegisterHandler(
		bot.HandlerTypeMessageText,
		"/start",
		bot.MatchTypeExact,
		startHandler.Start,
	)

	b.RegisterHandler(
		bot.HandlerTypeMessageText,
		"/menu",
		bot.MatchTypeExact,
		mainmenu.MainMenuHandler,
	)

	b.RegisterHandler(
		bot.HandlerTypeMessageText,
		"/info",
		bot.MatchTypeExact,
		mainmenu.InfoHandler,
	)

	expHandler := expense.NewHandler(expenseService)
	
	b.RegisterHandler(
		bot.HandlerTypeMessageText,
		"/expenses",
		bot.MatchTypeExact,
		expHandler.GetExpensesHandler,
	)

	b.RegisterHandler(
		bot.HandlerTypeMessageText,
		"/expense",
		bot.MatchTypePrefix,
		expHandler.AddExpenseHandler,
	)

	vidHandler := videohandler.NewHandler(videoService)
	b.RegisterHandlerMatchFunc(
    bot.MatchFunc(func(update *models.Update) bool { return true }),
    vidHandler.HandleMessage,
)
}

// func DefaultHandler(ctx context.Context, b *bot.Bot, update *models.Update) {
// 	if update.Message == nil || update.Message.Text == "" {
// 		return
// 	}

// 	b.SendMessage(ctx, &bot.SendMessageParams{
// 		ChatID: update.Message.Chat.ID,
// 		Text:   "Invalid command, type /info for more information and list of commands",
// 	})

// 	slog.Info(
// 		"defaultHandler message info: \n",
// 		slog.Int64("chat_id", update.Message.Chat.ID),
// 		slog.Int64("user_id", update.Message.From.ID),
// 		slog.String("user_name", update.Message.From.Username),
// 		slog.String("text", update.Message.Text),
// 	)
// }
