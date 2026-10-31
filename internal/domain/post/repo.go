package post

import (
	"context"
	"database/sql"
	"errors"

	"github.com/clary-work01/hot-post-api/internal/domain"
	"github.com/google/uuid"
	"gorm.io/gorm"
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

	if errors.Is(err, gorm.ErrDuplicatedKey) {
		return domain.ErrConflict
	}

	return err
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
		return nil, err
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
		return nil, 0, err
	}

	return list, total, nil
}
