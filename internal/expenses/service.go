package expenses

import (
	"BotinokTG/internal/chats"
	"BotinokTG/internal/members"
	"BotinokTG/internal/users"
	"context"
	"errors"
	"math"
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

func (s *ExpenseService) CreateEvent(
	ctx context.Context,
	tgChatID int64,
	title string,
	taggedUsernames []string,
	creatorTgID int64,
) (*Event, []string, error) {
	chat, err := s.chatRepo.GetByTelegramID(ctx, tgChatID)
	if err != nil {
		return nil, nil, errors.New("чат не найден, используйте /start")
	}

	chatMembers, err := s.memberRepo.GetMembersByChat(ctx, tgChatID)
	if err != nil {
		return nil, nil, err
	}

	memberByUsername := make(map[string]members.Member)
	memberByID := make(map[int64]members.Member)
	for _, m := range chatMembers {
		if m.Username.Valid {
			memberByUsername[m.Username.String] = m
		}
		memberByID[m.ID] = m
	}

	creator, err := s.userRepo.GetByTelegramID(ctx, creatorTgID)
	if err != nil {
		return nil, nil, errors.New("вы не зарегистрированы, используйте /start")
	}

	seen := make(map[int64]bool)
	var userIDs []int64
	var participantNames []string

	for _, tag := range taggedUsernames {
		m, ok := memberByUsername[tag]
		if !ok {
			continue
		}
		if !seen[m.ID] {
			seen[m.ID] = true
			userIDs = append(userIDs, m.ID)
			name := m.FirstName
			if m.Username.Valid {
				name = "@" + m.Username.String
			}
			participantNames = append(participantNames, name)
		}
	}

	if len(userIDs) == 0 {
		seen[creator.ID] = true
		userIDs = append(userIDs, creator.ID)
		name := creator.FirstName
		if cm, ok := memberByID[creator.ID]; ok && cm.Username.Valid {
			name = "@" + cm.Username.String
		}
		participantNames = append(participantNames, name)
	}

	event, err := s.expenseRepo.CreateEvent(ctx, chat.ID, title, userIDs)
	if err != nil {
		return nil, nil, err
	}

	return event, participantNames, nil
}

type EventSummary struct {
	Event            Event
	ParticipantCount int
	TotalAmount      int64
}

func (s *ExpenseService) GetActiveEventSummaries(
	ctx context.Context,
	tgChatID int64,
) ([]EventSummary, error) {
	return s.getEventSummaries(ctx, tgChatID, "open")
}

func (s *ExpenseService) GetClosedEventSummaries(
	ctx context.Context,
	tgChatID int64,
) ([]EventSummary, error) {
	return s.getEventSummaries(ctx, tgChatID, "closed")
}

func (s *ExpenseService) getEventSummaries(
	ctx context.Context,
	tgChatID int64,
	status string,
) ([]EventSummary, error) {
	chat, err := s.chatRepo.GetByTelegramID(ctx, tgChatID)
	if err != nil {
		return nil, errors.New("чат не найден")
	}

	var events []Event
	if status == "open" {
		events, err = s.expenseRepo.GetActiveEvents(ctx, chat.ID)
	} else {
		events, err = s.expenseRepo.GetClosedEvents(ctx, chat.ID)
	}
	if err != nil {
		return nil, err
	}

	var summaries []EventSummary
	for _, e := range events {
		parts, _ := s.expenseRepo.GetEventParticipants(ctx, e.ID)
		exps, _ := s.expenseRepo.GetEventExpenses(ctx, e.ID)
		var total int64
		for _, exp := range exps {
			total += exp.Amount
		}
		summaries = append(summaries, EventSummary{
			Event:            e,
			ParticipantCount: len(parts),
			TotalAmount:      total,
		})
	}

	return summaries, nil
}

type ExpenseRow struct {
	Expense   EventExpense
	PayerName string
}

type EventInfo struct {
	Event            Event
	ParticipantNames []string
	Expenses         []ExpenseRow
	TotalAmount      int64
}

func (s *ExpenseService) GetEventInfo(ctx context.Context, eventID int64) (*EventInfo, error) {
	ev, err := s.expenseRepo.GetEvent(ctx, eventID)
	if err != nil {
		return nil, errors.New("счет не найден")
	}

	parts, err := s.expenseRepo.GetEventParticipants(ctx, eventID)
	if err != nil {
		return nil, err
	}

	exps, err := s.expenseRepo.GetEventExpenses(ctx, eventID)
	if err != nil {
		return nil, err
	}

	userNames := make(map[int64]string)
	var partNames []string
	for _, p := range parts {
		userNames[p.UserID] = p.FirstName
		partNames = append(partNames, p.FirstName)
	}

	var rows []ExpenseRow
	var total int64
	for _, e := range exps {
		total += e.Amount
		name := userNames[e.PayerID]
		if name == "" {
			name = "Неизвестный"
		}
		rows = append(rows, ExpenseRow{Expense: e, PayerName: name})
	}

	return &EventInfo{
		Event:            *ev,
		ParticipantNames: partNames,
		Expenses:         rows,
		TotalAmount:      total,
	}, nil
}

func (s *ExpenseService) AddExpenseToEvent(
	ctx context.Context,
	eventID, tgPayerID, amount int64,
	description string,
) (string, error) {
	payer, err := s.userRepo.GetByTelegramID(ctx, tgPayerID)
	if err != nil {
		return "", errors.New("пользователь не найден, используйте /start")
	}

	ev, err := s.expenseRepo.GetEvent(ctx, eventID)
	if err != nil {
		return "", errors.New("счет не найден")
	}
	if ev.Status != "open" {
		return "", errors.New("этот счет уже закрыт")
	}

	parts, err := s.expenseRepo.GetEventParticipants(ctx, eventID)
	if err != nil {
		return "", err
	}

	isPart := false
	for _, p := range parts {
		if p.UserID == payer.ID {
			isPart = true

			break
		}
	}
	if !isPart {
		return "", errors.New("вы не участник этого счета")
	}

	_, err = s.expenseRepo.AddExpense(ctx, eventID, payer.ID, amount, description)

	return ev.Title, err
}

func (s *ExpenseService) CloseEvent(ctx context.Context, eventID int64) (string, error) {
	ev, err := s.expenseRepo.GetEvent(ctx, eventID)
	if err != nil {
		return "", errors.New("счет не найден")
	}

	return ev.Title, s.expenseRepo.CloseEvent(ctx, eventID)
}

type DebtResult struct {
	EventTitle       string
	TotalAmount      int64
	ParticipantCount int
	SharePerPerson   int64
	Debts            []Debt
}

func (s *ExpenseService) CalculateDebts(ctx context.Context, eventID int64) (*DebtResult, error) {
	ev, err := s.expenseRepo.GetEvent(ctx, eventID)
	if err != nil {
		return nil, errors.New("счет не найден")
	}

	parts, err := s.expenseRepo.GetEventParticipants(ctx, eventID)
	if err != nil || len(parts) == 0 {
		return nil, err
	}

	exps, err := s.expenseRepo.GetEventExpenses(ctx, eventID)
	if err != nil {
		return nil, err
	}

	var totalAmount int64
	paidByUser := make(map[int64]int64)
	for _, exp := range exps {
		totalAmount += exp.Amount
		paidByUser[exp.PayerID] += exp.Amount
	}

	share := totalAmount / int64(len(parts))

	for i := range parts {
		parts[i].Balance = paidByUser[parts[i].UserID] - share
	}

	var debtors []*UserBalance
	var creditors []*UserBalance

	for i := range parts {
		if parts[i].Balance < 0 {
			debtors = append(debtors, &parts[i])
		} else if parts[i].Balance > 0 {
			creditors = append(creditors, &parts[i])
		}
	}

	var debts []Debt
	dIdx, cIdx := 0, 0

	for dIdx < len(debtors) && cIdx < len(creditors) {
		debtor := debtors[dIdx]
		creditor := creditors[cIdx]

		owe := int64(math.Abs(float64(debtor.Balance)))
		receive := creditor.Balance

		amount := owe
		if receive < owe {
			amount = receive
		}

		if amount > 0 {
			debts = append(debts, Debt{
				DebtorID:     debtor.UserID,
				DebtorName:   debtor.FirstName,
				CreditorID:   creditor.UserID,
				CreditorName: creditor.FirstName,
				Amount:       amount,
			})
		}

		debtor.Balance += amount
		creditor.Balance -= amount

		if debtor.Balance == 0 {
			dIdx++
		}
		if creditor.Balance == 0 {
			cIdx++
		}
	}

	return &DebtResult{
		EventTitle:       ev.Title,
		TotalAmount:      totalAmount,
		ParticipantCount: len(parts),
		SharePerPerson:   share,
		Debts:            debts,
	}, nil
}
