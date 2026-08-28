package who

import (
	"BotinokTG/internal/members"
	"context"
	"crypto/rand"
	"errors"
	"math/big"
)

type MemberProvider interface {
	GetMembersByChat(ctx context.Context, chatID int64) ([]members.Member, error)
}

type RandomizerService struct {
	repo MemberProvider
}

func NewRandomizerService(repo MemberProvider) *RandomizerService {
	return &RandomizerService{
		repo: repo,
	}
}

func (s *RandomizerService) GetRandomMember(
	ctx context.Context,
	chatID int64,
) (*members.Member, error) {
	membersList, err := s.repo.GetMembersByChat(ctx, chatID)
	if err != nil {
		return nil, err
	}

	if len(membersList) == 0 {
		return nil, errors.New("not enough members")
	}
	// крипто рандом, т.к. обычный будто бы не то
	n, err := rand.Int(rand.Reader, big.NewInt(int64(len(membersList))))
	if err != nil {
		return nil, err
	}

	randomIndex := n.Int64()

	return &membersList[randomIndex], nil
}
