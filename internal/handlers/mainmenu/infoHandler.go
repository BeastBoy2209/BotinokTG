package mainmenu

import (
	"context"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

func InfoHandler(ctx context.Context, b *bot.Bot, update *models.Update) {
	if update.Message == nil {
		return
	}

	msg := "👋 Привет! Я Botinok — помощник для групповых чатов.\n\n" +
		"💸 *Управление счетами:*\n" +
		"/event create [Название] [@user1] — создать счет\n" +
		"/event list — список активных счетов\n" +
		"/event history — история закрытых счетов\n" +
		"/event info [ID] — детали счета\n" +
		"/event add [ID] [Сумма] [Описание] — добавить расход\n" +
		"/event debts [ID] — рассчитать долги\n" +
		"/event close [ID] — закрыть счет\n\n" +
		"🎥 *Скачивание видео:*\n" +
		"Просто отправь ссылку на TikTok, YouTube Shorts, Instagram Reels или Pinterest, и я пришлю видео без водяных знаков.\n\n" +
		"🎲 *Случайный выбор и теги:*\n" +
		"/who — выбрать случайного участника чата\n" +
		"/who @user1 @user2 — выбрать случайного из списка\n" +
		"@all или @все — тегнуть всех активных участников группы"

	_, _ = b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID:    update.Message.Chat.ID,
		Text:      msg,
		ParseMode: models.ParseModeMarkdown,
	})
}
