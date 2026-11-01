package auth

import (
	"errors"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

var (
	ErrInvalidRefreshToken = errors.New("invalid refresh token")
	ErrInvalidCredentials  = errors.New("invalid username or password")
)

type RefreshToken struct {
	ID        uuid.UUID
	UserID    uuid.UUID
	TokenHash string
	ExpiresAt time.Time
	CreatedAt time.Time
	RevokedAt *time.Time
}

func (u *RefreshToken) BeforeCreate(tx *gorm.DB) error {
	if u.ID == uuid.Nil {
		u.ID = uuid.Must(uuid.NewV7())
	}
	return nil
}
