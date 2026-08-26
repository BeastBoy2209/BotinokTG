package chats

import (
	"time"

	"github.com/jackc/pgx/v5/pgtype"
)

type Chat struct {
	ID         int64
	TelegramID int64
	Title      pgtype.Text
	CreatedAt  time.Time
}
