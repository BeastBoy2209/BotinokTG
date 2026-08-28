package expenses

import (
	"BotinokTG/internal/storage"
	"context"

	"github.com/jackc/pgx/v5/pgtype"
)

type ExpenseRepository struct {
	db storage.DBTX
}

func NewRepository(db storage.DBTX) *ExpenseRepository {
	return &ExpenseRepository{db: db}
}

func (r *ExpenseRepository) CreateEvent(
	ctx context.Context,
	chatID int64,
	title string,
	userIDs []int64,
) (*Event, error) {
	var event Event

	query := `
	INSERT INTO events (chat_id, title, status)
	VALUES ($1, $2, 'open')
	RETURNING id, chat_id, title, status, created_at
	`
	err := r.db.QueryRow(ctx, query, chatID, title).Scan(
		&event.ID,
		&event.ChatID,
		&event.Title,
		&event.Status,
		&event.CreatedAt,
	)
	if err != nil {
		return nil, err
	}

	for _, userID := range userIDs {
		partQuery := `INSERT INTO event_participants (event_id, user_id) VALUES ($1, $2) ON CONFLICT DO NOTHING`
		_, err := r.db.Exec(ctx, partQuery, event.ID, userID)
		if err != nil {
			return nil, err
		}
	}

	return &event, nil
}

func (r *ExpenseRepository) GetActiveEvents(ctx context.Context, chatID int64) ([]Event, error) {
	return r.getEventsByStatus(ctx, chatID, "open")
}

func (r *ExpenseRepository) GetClosedEvents(ctx context.Context, chatID int64) ([]Event, error) {
	return r.getEventsByStatus(ctx, chatID, "closed")
}

func (r *ExpenseRepository) getEventsByStatus(
	ctx context.Context,
	chatID int64,
	status string,
) ([]Event, error) {
	query := `SELECT id, chat_id, title, status, created_at FROM events WHERE chat_id = $1 AND status = $2 ORDER BY created_at DESC`
	rows, err := r.db.Query(ctx, query, chatID, status)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var events []Event
	for rows.Next() {
		var e Event
		err := rows.Scan(&e.ID, &e.ChatID, &e.Title, &e.Status, &e.CreatedAt)
		if err != nil {
			return nil, err
		}
		events = append(events, e)
	}

	return events, nil
}

func (r *ExpenseRepository) GetEvent(ctx context.Context, eventID int64) (*Event, error) {
	var e Event
	query := `SELECT id, chat_id, title, status, created_at FROM events WHERE id = $1`
	err := r.db.QueryRow(ctx, query, eventID).
		Scan(&e.ID, &e.ChatID, &e.Title, &e.Status, &e.CreatedAt)
	if err != nil {
		return nil, err
	}

	return &e, nil
}

func (r *ExpenseRepository) CloseEvent(ctx context.Context, eventID int64) error {
	query := `UPDATE events SET status = 'closed' WHERE id = $1`
	_, err := r.db.Exec(ctx, query, eventID)

	return err
}

func (r *ExpenseRepository) AddExpense(
	ctx context.Context,
	eventID, payerID, amount int64,
	description string,
) (*EventExpense, error) {
	var expense EventExpense
	var pgDescription pgtype.Text

	if description != "" {
		pgDescription = pgtype.Text{String: description, Valid: true}
	} else {
		pgDescription = pgtype.Text{Valid: false}
	}

	query := `
	INSERT INTO event_expenses (event_id, payer_id, amount, description) 
	VALUES ($1, $2, $3, $4) 
	RETURNING id, event_id, payer_id, amount, description, created_at
	`
	err := r.db.QueryRow(ctx, query, eventID, payerID, amount, pgDescription).Scan(
		&expense.ID,
		&expense.EventID,
		&expense.PayerID,
		&expense.Amount,
		&expense.Description,
		&expense.CreatedAt,
	)
	if err != nil {
		return nil, err
	}

	return &expense, nil
}

func (r *ExpenseRepository) GetEventParticipants(
	ctx context.Context,
	eventID int64,
) ([]UserBalance, error) {
	query := `
	SELECT u.id, u.first_name
	FROM users u
	JOIN event_participants ep ON ep.user_id = u.id
	WHERE ep.event_id = $1
	`
	rows, err := r.db.Query(ctx, query, eventID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var balances []UserBalance
	for rows.Next() {
		var ub UserBalance
		err := rows.Scan(&ub.UserID, &ub.FirstName)
		if err != nil {
			return nil, err
		}
		balances = append(balances, ub)
	}

	return balances, nil
}

func (r *ExpenseRepository) GetEventExpenses(
	ctx context.Context,
	eventID int64,
) ([]EventExpense, error) {
	query := `SELECT id, event_id, payer_id, amount, description, created_at FROM event_expenses WHERE event_id = $1`
	rows, err := r.db.Query(ctx, query, eventID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var exps []EventExpense
	for rows.Next() {
		var e EventExpense
		err := rows.Scan(
			&e.ID,
			&e.EventID,
			&e.PayerID,
			&e.Amount,
			&e.Description,
			&e.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		exps = append(exps, e)
	}

	return exps, nil
}
