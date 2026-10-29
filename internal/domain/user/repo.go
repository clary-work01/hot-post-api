package user

import (
	"context"
	"errors"

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

func (r *Repo) follow(ctx context.Context, follow *Follow) (inserted bool, err error) {
	res := r.db.WithContext(ctx).
		Clauses(clause.OnConflict{DoNothing: true}).
		Create(follow)

	err = res.Error
	// gorm.ErrForeignKeyViolated 只有在開啟 TranslateError: true 時才會出現
	if errors.Is(err, gorm.ErrForeignKeyViolated) {
		return false, domain.ErrNotFound
	}

	if err != nil {
		return false, err
	}

	return res.RowsAffected == 1, nil
}

func (r *Repo) unFollow(ctx context.Context, followerID, followeeID uuid.UUID) (inserted bool, err error) {
	res := r.db.WithContext(ctx).
		Where("follower_id = ? AND followee_id = ?", followerID, followeeID).
		Delete(&Follow{})

	err = res.Error
	// gorm.ErrForeignKeyViolated 只有在開啟 TranslateError: true 時才會出現
	if errors.Is(err, gorm.ErrForeignKeyViolated) {
		return false, domain.ErrNotFound
	}

	if err != nil {
		return false, err
	}

	return res.RowsAffected == 1, nil
}
