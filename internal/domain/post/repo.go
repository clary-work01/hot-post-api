package post

import (
	"context"
	"database/sql"
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

func (r *Repo) create(ctx context.Context, post *Post) error {
	err := r.db.WithContext(ctx).
		Create(post).Error

	if err != nil {
		return fmt.Errorf("create post (author_id %s): %w", post.AuthorID, err)
	}

	return nil
}

const isLikedExpr = `EXISTS (
	SELECT 1 FROM post_likes pl
	WHERE pl.user_id = @viewerID AND pl.post_id = p.id
) AS is_liked`

func (r *Repo) detail(ctx context.Context, postID, viewerID uuid.UUID) (*Detail, error) {
	var detail Detail

	err := r.db.WithContext(ctx).
		Table("posts AS p").
		Select(`
			p.id,
			p.author_id, 
			u.username AS author_name,
			p.title,
			p.content,
			p.like_count,
			p.created_at,
			p.updated_at,
			`+isLikedExpr, sql.Named("viewerID", viewerID)).
		Joins("JOIN users u ON u.id = p.author_id").
		Where("p.id = ?", postID).
		Take(&detail).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, domain.ErrNotFound
	}

	if err != nil {
		return nil, fmt.Errorf("get post detail (%s): %w", postID, err)
	}

	return &detail, nil
}

const excerptLen = 50

func (r *Repo) list(ctx context.Context, offset, limit int, viewerID uuid.UUID) (list []Summary, total int64, err error) {
	list = make([]Summary, 0, limit)

	err = r.db.WithContext(ctx).
		Table("posts AS p").
		Joins("JOIN users u ON u.id = p.author_id").
		Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	err = r.db.WithContext(ctx).
		Table("posts AS p").
		Select(`
			p.id, 
			p.author_id, 
			u.username AS author_name,
			p.title,
			LEFT(p.content, @excerptLen) AS excerpt,
			p.like_count,
			p.created_at,
			p.updated_at,
			`+isLikedExpr, sql.Named("excerptLen", excerptLen), sql.Named("viewerID", viewerID)).
		Joins("JOIN users u ON u.id = p.author_id").
		Order("p.created_at DESC, p.id DESC").
		Offset(offset).
		Limit(limit).
		Find(&list).Error

	if err != nil {
		return nil, 0, fmt.Errorf("list posts: %w", err)
	}

	return list, total, nil
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
