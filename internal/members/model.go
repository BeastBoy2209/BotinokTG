package members

import (
	"time"
)

type ChatMember struct {
	ChatID   int64
	UserID   int64
	JoinedAt time.Time
}
