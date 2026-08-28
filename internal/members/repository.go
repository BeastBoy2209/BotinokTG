package members

import (
	"BotinokTG/internal/storage"
	"context"
	"time"

	"github.com/jackc/pgx/v5/pgtype"
)

type MembershipRepository struct {
	db storage.DBTX
}

func NewRepository(db storage.DBTX) *MembershipRepository {
	return &MembershipRepository{
		db: db,
	}
}

func (r *MembershipRepository) AddMember(ctx context.Context, chatID, userID int64) error {
	query := `
		INSERT INTO chat_members (chat_id, user_id)
		VALUES ($1, $2)
		ON CONFLICT (chat_id, user_id) DO NOTHING
	`
	_, err := r.db.Exec(ctx, query, chatID, userID)

	return err
}

func (r *MembershipRepository) IsMember(ctx context.Context, chatID, userID int64) (bool, error) {
	var isMember bool
	query := `
		SELECT EXISTS(
		SELECT 1 FROM chat_members 
		WHERE chat_id = $1 AND user_id = $2)
	`

	err := r.db.QueryRow(ctx, query, chatID, userID).Scan(&isMember)
	if err != nil {
		return false, err
	}

	return isMember, nil
}

type Member struct {
	ID         int64
	TelegramID int64
	Username   pgtype.Text
	FirstName  string
	CreatedAt  time.Time
}

func (r *MembershipRepository) GetMembersByChat(
	ctx context.Context,
	chatID int64,
) ([]Member, error) {
	query := `
		SELECT u.id, u.telegram_id, u.username, u.first_name, u.created_at
		FROM users u
		JOIN chat_members cm ON u.id = cm.user_id
		JOIN chats c ON c.id = cm.chat_id
		WHERE c.telegram_id = $1
	`
	rows, err := r.db.Query(ctx, query, chatID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var chatMembers []Member
	for rows.Next() {
		var user Member
		err := rows.Scan(
			&user.ID,
			&user.TelegramID,
			&user.Username,
			&user.FirstName,
			&user.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		chatMembers = append(chatMembers, user)
	}

	return chatMembers, rows.Err()
}
