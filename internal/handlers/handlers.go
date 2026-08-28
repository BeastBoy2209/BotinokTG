package handlers

import (
	"BotinokTG/internal/expenses"
	"BotinokTG/internal/handlers/expense"
	"BotinokTG/internal/handlers/mainmenu"
	"BotinokTG/internal/handlers/start"
	"BotinokTG/internal/handlers/videohandler"
	"BotinokTG/internal/members"
	"BotinokTG/internal/users"
	"BotinokTG/internal/video"
	"BotinokTG/internal/who"
	"context"

	// "log/slog".

	mentionhandler "BotinokTG/internal/handlers/mention"
	whohandler "BotinokTG/internal/handlers/who"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

func RegistrationMiddleware(service *users.RegistrationService) bot.Middleware {
	return func(next bot.HandlerFunc) bot.HandlerFunc {
		return func(ctx context.Context, b *bot.Bot, update *models.Update) {
			if update.Message != nil && update.Message.From != nil {
				var chatTitle string
				if update.Message.Chat.Type == "private" {
					chatTitle = update.Message.Chat.FirstName
				} else {
					chatTitle = update.Message.Chat.Title
				}

				go func() {
					_ = service.RegisterUserInChat(
						context.Background(),
						update.Message.From.ID,
						update.Message.From.Username,
						update.Message.From.FirstName,
						update.Message.Chat.ID,
						chatTitle,
					)
				}()
			}
			next(ctx, b, update)
		}
	}
}

func RegisterAll(
	b *bot.Bot,
	startService *users.RegistrationService,
	expenseService *expenses.ExpenseService,
	videoService *video.VideoService,
	whoService *who.RandomizerService,
	memberRepo *members.MembershipRepository,
) {
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
		"/event",
		bot.MatchTypePrefix,
		expHandler.HandleEventCommand,
	)

	whoHand := whohandler.NewHandler(whoService)
	b.RegisterHandler(
		bot.HandlerTypeMessageText,
		"/who",
		bot.MatchTypePrefix,
		whoHand.HandleWho,
	)

	mentionHand := mentionhandler.NewHandler(memberRepo)
	b.RegisterHandlerMatchFunc(
		bot.MatchFunc(mentionhandler.MatchFunc),
		mentionHand.Handle,
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
