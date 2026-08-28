package mention

import (
	"BotinokTG/internal/members"
	"context"
	"log/slog"
	"strings"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

var triggers = []string{"@all", "@все", "!all", "!все", "/all", "/все"}

type Handler struct {
	repo *members.MembershipRepository
}

func NewHandler(repo *members.MembershipRepository) *Handler {
	return &Handler{repo: repo}
}

func MatchFunc(update *models.Update) bool {
	if update.Message == nil || update.Message.Text == "" {
		return false
	}
	lower := strings.ToLower(update.Message.Text)
	for _, t := range triggers {
		if strings.Contains(lower, t) {
			return true
		}
	}

	return false
}

func (h *Handler) Handle(ctx context.Context, b *bot.Bot, update *models.Update) {
	chatID := update.Message.Chat.ID

	chatMembers, err := h.repo.GetMembersByChat(ctx, chatID)
	if err != nil {
		slog.Error("mention all: failed to get members", slog.Any("error", err))

		return
	}
	if len(chatMembers) == 0 {
		return
	}

	text := "👋 "
	var entities []models.MessageEntity

	for _, m := range chatMembers {
		var part string
		if m.Username.Valid && m.Username.String != "" {
			part = "@" + m.Username.String
		} else {
			part = m.FirstName

			offset := utf16Len(text)
			entities = append(entities, models.MessageEntity{
				Type:   models.MessageEntityTypeTextMention,
				Offset: offset,
				Length: utf16Len(part),
				User: &models.User{
					ID: m.TelegramID,
				},
			})
		}

		text += part + " "
	}

	params := &bot.SendMessageParams{
		ChatID: chatID,
		Text:   text,
	}
	if len(entities) > 0 {
		params.Entities = entities
	}

	_, err = b.SendMessage(ctx, params)
	if err != nil {
		slog.Error("mention all: failed to send", slog.Any("error", err))
	}
}

func utf16Len(s string) int {
	n := 0
	for _, r := range s {
		if r <= 0xFFFF {
			n++
		} else {
			n += 2
		}
	}

	return n
}
