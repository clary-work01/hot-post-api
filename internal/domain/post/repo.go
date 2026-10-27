package post

import (
	"context"
	"errors"

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
