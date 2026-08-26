package chats

import (
	"context"

	"BotinokTG/internal/storage"

	"github.com/jackc/pgx/v5/pgtype"
)

type ChatRepository struct {
	db storage.DBTX
}

func NewRepository(db storage.DBTX) *ChatRepository {
	return &ChatRepository{
		db: db,
	}
}

func (r *ChatRepository) UpsertChat(ctx context.Context, telegramID int64, title string) (*Chat, error) {
	var chat Chat

	var pgTitle pgtype.Text
	if title != "" {
		pgTitle = pgtype.Text{String: title, Valid: true}
	} else {
		pgTitle = pgtype.Text{Valid: false}
	}

	query := `
		INSERT INTO chats (telegram_id, title)
		VALUES ($1, $2)
		ON CONFLICT (telegram_id) DO UPDATE
		SET title = EXCLUDED.title
		RETURNING id, telegram_id, title, created_at
	`

	err := r.db.QueryRow(ctx, query, telegramID, pgTitle).Scan(
		&chat.ID,
		&chat.TelegramID,
		&chat.Title,
		&chat.CreatedAt,
	)
	if err != nil {
		return nil, err
	}

	return &chat, nil
}

func (r *ChatRepository) GetByTelegramID(ctx context.Context, telegramID int64) (*Chat, error) {
	var chat Chat
	query := `
	SELECT * FROM chats WHERE telegram_id = $1
	`

	err := r.db.QueryRow(ctx, query, telegramID).Scan(
		&chat.ID,
		&chat.TelegramID,
		&chat.Title,
		&chat.CreatedAt,
	)
	if err != nil {
		return nil, err
	}

	return &chat, nil
}
