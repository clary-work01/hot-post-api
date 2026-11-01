package auth

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/google/uuid"
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

// 作廢舊的refresh token ＋ 寫入新的，同一個 transaction
func (r *Repo) rotateRefreshToken(ctx context.Context, oldTokenHash string, newToken *RefreshToken) (uuid.UUID, error) {
	var userID uuid.UUID

	err := r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {

		// 同時「驗證＋作廢」
		err := tx.Raw(`
			UPDATE refresh_tokens
			SET revoked_at = NOW()
			WHERE token_hash = ? AND revoked_at IS NULL AND expires_at > NOW()
			RETURNING user_id
		`, oldTokenHash).Row().Scan(&userID)

		if errors.Is(err, sql.ErrNoRows) {
			return ErrInvalidRefreshToken // 不存在 / 已用過 / 過期
		}
		if err != nil {
			return fmt.Errorf("revoke old refresh token: %w", err)
		}

		// 寫入新的 refresh token
		newToken.UserID = userID
		if err := tx.Create(newToken).Error; err != nil {
			return fmt.Errorf("create new refresh token: %w", err)
		}

		return nil
	})

	if err != nil {
		return uuid.Nil, fmt.Errorf("rotate refresh token: %w", err)
	}

	return userID, nil
}

func (r *Repo) revokeRefreshToken(ctx context.Context, tokenHash string) error {
	err := r.db.WithContext(ctx).
		Model(&RefreshToken{}).
		Where("token_hash = ? AND revoked_at IS NULL", tokenHash).
		Update("revoked_at", gorm.Expr("NOW()")).Error

	if err != nil {
		return fmt.Errorf("revoke refresh token: %w", err)
	}

	return nil
}
