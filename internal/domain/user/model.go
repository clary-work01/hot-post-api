package user

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID        uuid.UUID `gorm:"default:uuidv7()"`
	Username  string
	Password  string
	CreatedAt time.Time
}

type Follow struct {
	FollowerID uuid.UUID
	FolloweeID uuid.UUID
	CreatedAt  time.Time
}

type Profile struct {
	ID             uuid.UUID
	Username       string
	PostCount      int
	FollowerCount  int
	FollowingCount int
	IsFollowing    bool
}

// FollowItem: 粉絲/關注列表中的一筆
type FollowItem struct {
	ID          uuid.UUID
	Username    string
	FollowedAt  time.Time
	IsFollowing bool
}
