package user

import (
	"errors"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

var (
	ErrCannotFollowSelf = errors.New("cannot follow yourself")
	ErrUsernameTaken    = errors.New("user name is already taken")
)

type User struct {
	ID        uuid.UUID
	Username  string
	Password  string
	CreatedAt time.Time
}

func (u *User) BeforeCreate(tx *gorm.DB) error {
	if u.ID == uuid.Nil {
		u.ID = uuid.Must(uuid.NewV7())
	}
	return nil
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
