package expenses

import (
	"BotinokTG/internal/chats"
	"BotinokTG/internal/members"
	"BotinokTG/internal/users"
	"context"
	"errors"
)

type ExpenseService struct {
	expenseRepo *ExpenseRepository
	memberRepo  *members.MembershipRepository
	userRepo    *users.UserRepository
	chatRepo    *chats.ChatRepository
}

func NewService(
	expenseRepo *ExpenseRepository,
	memberRepo *members.MembershipRepository,
	userRepo *users.UserRepository,
	chatRepo *chats.ChatRepository,
) *ExpenseService {
	return &ExpenseService{
		expenseRepo: expenseRepo,
		memberRepo:  memberRepo,
		userRepo:    userRepo,
		chatRepo:    chatRepo,
	}
}

func (s *ExpenseService) AddExpense(ctx context.Context, tgChatID, tgUserID, amount int64, description string) (*Expense, error){
	user, err := s.userRepo.GetByTelegramID(ctx, tgUserID)
	if err != nil{
		return nil, errors.New("вы не зарегистрированы, нажмите /start")
	}
	chat, err := s.chatRepo.GetByTelegramID(ctx, tgChatID)
	if err != nil{
		return nil, errors.New("бот не добавлен в этот чат, нажмите /start")
	}
	isMember, err := s.memberRepo.IsMember(ctx, chat.ID, user.ID)
	if err != nil{
		return nil, errors.New("бот не добавлен в этот чат, нажмите /start")
	}
	if !isMember{
		return nil, errors.New("пользователь не является участником чата")
	}

	return s.expenseRepo.CreateExpense(ctx, chat.ID, user.ID, amount, description)
}

func (s *ExpenseService) GetChatExpenses(ctx context.Context, tgChatID int64) ([]ChatExpenseInfo, error){
	chat, err := s.chatRepo.GetByTelegramID(ctx, tgChatID)
	if err != nil{
		return nil, errors.New("бот не добавлен в этот чат, нажмите /start")
	}
	
	return s.expenseRepo.GetExpensesByChat(ctx, chat.ID)
}
