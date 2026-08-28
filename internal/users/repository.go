package users

import (
	"BotinokTG/internal/storage"
	"context"

	"github.com/jackc/pgx/v5/pgtype"
)

type UserRepository struct {
	db storage.DBTX
}

func NewRepository(db storage.DBTX) *UserRepository {
	return &UserRepository{
		db: db,
	}
}

func (r *UserRepository) UpsertUser(
	ctx context.Context,
	telegram_id int64,
	username, firstname string,
) (*User, error) {
	var user User

	var pgUsername pgtype.Text
	if username != "" {
		pgUsername = pgtype.Text{String: username, Valid: true}
	} else {
		pgUsername = pgtype.Text{Valid: false}
	}

	query := `
    INSERT INTO users (telegram_id, username, first_name) 
    VALUES ($1, $2, $3) 
    ON CONFLICT (telegram_id) DO UPDATE 
    SET username = EXCLUDED.username, first_name = EXCLUDED.first_name 
    RETURNING id, telegram_id, username, first_name, created_at
	`
	err := r.db.QueryRow(ctx, query, telegram_id, pgUsername, firstname).Scan(
		&user.ID,
		&user.TelegramID,
		&user.Username,
		&user.FirstName,
		&user.CreatedAt,
	)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *UserRepository) GetByTelegramID(ctx context.Context, telegramID int64) (*User, error) {
	var user User
	query := `
		SELECT id, telegram_id, username, first_name, created_at
		FROM users 
		WHERE telegram_id = $1
	`

	err := r.db.QueryRow(ctx, query, telegramID).Scan(
		&user.ID,
		&user.TelegramID,
		&user.Username,
		&user.FirstName,
		&user.CreatedAt,
	)
	if err != nil {
		return nil, err
	}

	return &user, nil
}
