package members

import (
	"context"

	"BotinokTG/internal/storage"
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

func (r *MembershipRepository) IsMember(ctx context.Context, chatID int64, userID int64) (bool, error) {
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
