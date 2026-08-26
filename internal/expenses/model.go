package expenses

import (
	"time"

	"github.com/jackc/pgx/v5/pgtype"
)

type Expense struct {
	ID          int64
	ChatID      int64
	UserID      int64
	Amount      int64
	Description pgtype.Text
	CreatedAt   time.Time
}

type ChatExpenseInfo struct {
	Amount      int64
	Description pgtype.Text
	CreatedAt   time.Time
	Username    pgtype.Text
	FirstName   pgtype.Text
}
