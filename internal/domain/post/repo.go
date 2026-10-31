package post

import (
	"context"
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
		return nil
	}

	return err
}

func (r *Repo) detail(ctx context.Context, postID, viewerID uuid.UUID) (*PostDetail, error) {
	var PostDetail PostDetail

	err := r.db.WithContext(ctx).
		Table("posts AS p").
		Select(`
			p.id,p.author_id,u.username AS author_name,p.title,p.content,p.like_count,p.created_at,p.updated_at,
			EXISTS (
				SELECT 1 FROM post_likes pl
				WHERE pl.user_id = ? AND pl.post_id = p.id
			) AS is_liked
			`, viewerID).
		Joins("JOIN users u ON u.id = p.author_id").
		Where("p.id = ?", postID).
		Take(&PostDetail).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, domain.ErrNotFound
	}

	if err != nil {
		return nil, err
	}

	return &PostDetail, nil
}
