package expenses

import (
	"context"

	"BotinokTG/internal/storage"

	"github.com/jackc/pgx/v5/pgtype"
)

type ExpenseRepository struct {
	db storage.DBTX
}

func NewRepository(db storage.DBTX) *ExpenseRepository {
	return &ExpenseRepository{db: db}
}

func (r *ExpenseRepository) CreateExpense(ctx context.Context, chatID int64, userID int64, amount int64, description string) (*Expense, error) {
	var expense Expense
	var pgDescription pgtype.Text

	if description != "" {
		pgDescription = pgtype.Text{String: description, Valid: true}
	} else {
		pgDescription = pgtype.Text{Valid: false}
	}

	query := `
	INSERT INTO expenses (chat_id, user_id, amount, description) 
	VALUES ($1, $2, $3, $4) 
	RETURNING id, chat_id, user_id, amount, description, created_at
	`
	err := r.db.QueryRow(ctx, query, chatID, userID, amount, pgDescription).Scan(
		&expense.ID,
		&expense.ChatID,
		&expense.UserID,
		&expense.Amount,
		&expense.Description,
		&expense.CreatedAt,
	)
	if err != nil {
		return nil, err
	}

	return &expense, nil
}

func (r *ExpenseRepository) GetExpensesByChat(ctx context.Context, chatID int64) ([]ChatExpenseInfo, error) {
	query := `
	SELECT e.amount, e.description, e.created_at, u.username, u.first_name 
	FROM expenses e 
	JOIN users u 
	ON e.user_id = u.id 
	WHERE e.chat_id = $1 
	ORDER BY e.created_at DESC
	`

	rows, err := r.db.Query(ctx, query, chatID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	expenses := make([]ChatExpenseInfo, 0)
	for rows.Next() {
		var exp ChatExpenseInfo

		err := rows.Scan(
			&exp.Amount,
			&exp.Description,
			&exp.CreatedAt,
			&exp.Username,
			&exp.FirstName,
		)
		if err != nil {
			return nil, err
		}

		expenses = append(expenses, exp)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return expenses, nil
}
