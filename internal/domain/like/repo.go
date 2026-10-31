package like

import (
	"context"
	"errors"
	"fmt"

	"github.com/clary-work01/hot-post-api/internal/domain"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type Repo struct {
	db *gorm.DB
}

func NewRepo(db *gorm.DB) *Repo {
	return &Repo{db}
}

func (r *Repo) like(ctx context.Context, postLike *PostLike) (inserted bool, err error) {
	err = r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		res := tx.Clauses(clause.OnConflict{DoNothing: true}).
			Create(postLike)

		// 按讚的貼文不存在時，會回 FK violation
		if errors.Is(res.Error, gorm.ErrForeignKeyViolated) {
			return domain.ErrNotFound // 貼文不存在
		}

		if res.Error != nil {
			return fmt.Errorf("insert like: %w", res.Error)
		}

		// 已讚過
		if res.RowsAffected == 0 {
			return nil
		}

		err := tx.Table("posts").
			Where("id = ?", postLike.PostID).
			UpdateColumn("like_count", gorm.Expr("like_count + 1")).Error

		if err != nil {
			return fmt.Errorf("increment like count: %w", err)
		}

		inserted = true
		return nil
	})
	if errors.Is(err, domain.ErrNotFound) {
		return false, err
	}

	if err != nil {
		return false, fmt.Errorf("like post (post %s, user %s): %w", postLike.PostID, postLike.UserID, err)
	}

	return inserted, err
}

func (r *Repo) unlike(ctx context.Context, postID, userID uuid.UUID) (deleted bool, err error) {
	err = r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		res := tx.Where("post_id = ? AND user_id = ?", postID, userID).
			Delete(&PostLike{})

		if res.Error != nil {
			return fmt.Errorf("delete like: %w", res.Error)
		}

		// 本來就沒讚過（或貼文不存在）
		if res.RowsAffected == 0 {
			return nil
		}

		err := tx.Table("posts").
			Where("id = ?", postID).
			UpdateColumn("like_count", gorm.Expr("like_count -1")).Error

		if err != nil {
			return fmt.Errorf("decrement like count: %w", err)
		}

		deleted = true
		return nil
	})

	if err != nil {
		return false, fmt.Errorf("unlike post (post %s, user %s): %w", postID, userID, err)
	}

	return deleted, err
}
