package users

import (
	"BotinokTG/internal/chats"
	"BotinokTG/internal/members"
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

type RegistrationService struct {
	pool *pgxpool.Pool
}

func NewRegistrationService(pool *pgxpool.Pool) *RegistrationService {
	return &RegistrationService{
		pool: pool,
	}
}

func (s *RegistrationService) RegisterUserInChat(
	ctx context.Context,
	tgUserID int64,
	username, firstName string,
	tgChatID int64,
	chatTitle string,
) error {
	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	userRepo := NewRepository(tx)
	chatRepo := chats.NewRepository(tx)
	memberRepo := members.NewRepository(tx)

	user, err := userRepo.UpsertUser(ctx, tgUserID, username, firstName)
	if err != nil {
		return err
	}

	chat, err := chatRepo.UpsertChat(ctx, tgChatID, chatTitle)
	if err != nil {
		return err
	}

	err = memberRepo.AddMember(ctx, chat.ID, user.ID)
	if err != nil {
		return err
	}

	return tx.Commit(ctx)
}
