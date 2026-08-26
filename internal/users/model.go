package users

import (
	"github.com/jackc/pgx/v5/pgtype"
	"time"
)

type User struct {
	ID         int64
	TelegramID int64
	Username   pgtype.Text
	FirstName  string
	CreatedAt  time.Time
}
