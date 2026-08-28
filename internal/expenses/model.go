package expenses

import (
	"time"

	"github.com/jackc/pgx/v5/pgtype"
)

type Event struct {
	ID        int64
	ChatID    int64
	Title     string
	Status    string // 'open', 'closed'
	CreatedAt time.Time
}

type EventExpense struct {
	ID          int64
	EventID     int64
	PayerID     int64
	Amount      int64
	Description pgtype.Text
	CreatedAt   time.Time
}

type Debt struct {
	DebtorID     int64
	CreditorID   int64
	DebtorName   string
	CreditorName string
	Amount       int64
}

type UserBalance struct {
	UserID    int64
	FirstName string
	Balance   int64
}
