package auth

import (
	"context"
	"fmt"

	"gorm.io/gorm"
)

type Repo struct {
	db *gorm.DB
}

func NewRepo(db *gorm.DB) *Repo {
	return &Repo{db}
}

func (r *Repo) createRefreshToken(ctx context.Context, rt *RefreshToken) error {
	err := r.db.WithContext(ctx).Create(rt).Error

	if err != nil {
		return fmt.Errorf("create refresh token (user_id %s): %w", rt.UserID, err)
	}
	return nil
}
