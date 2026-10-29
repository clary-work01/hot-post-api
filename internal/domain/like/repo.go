package like

import (
	"context"

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

		if res.Error != nil {
			return res.Error
		}

		// 已讚過
		if res.RowsAffected == 0 {
			return nil
		}

		// 按讚成功 才更新posts.like_count
		inserted = true

		err := tx.Table("posts").
			Where("id = ?", postLike.PostID).
			UpdateColumn("like_count", gorm.Expr("like_count + ?", 1)).Error

		return err
	})

	// 按讚的貼文不存在時，會回 FK violation

	return inserted, err
}

func (r *Repo) unlike(ctx context.Context, postID, userID uuid.UUID) (deleted bool, err error) {
	err = r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		res := tx.Where("post_id = ? AND user_id = ?", postID, userID).
			Delete(&PostLike{})

		if res.Error != nil {
			return res.Error
		}

		// 本來就沒讚過
		if res.RowsAffected == 0 {
			return nil
		}

		// 刪除按讚成功 才更新posts.like_count
		deleted = true

		err := tx.Table("posts").
			Where("id = ?", postID).
			UpdateColumn("like_count", gorm.Expr("like_count + ?", -1)).Error

		return err
	})

	// 按讚的貼文不存在時，會回 FK violation

	return deleted, err
}
