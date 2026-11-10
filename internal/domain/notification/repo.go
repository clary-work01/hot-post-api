package notification

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

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

func (r *Repo) create(ctx context.Context, n *Notification) error {
	err := r.db.WithContext(ctx).
		Clauses(clause.OnConflict{DoNothing: true}).
		Create(n).Error
	if err != nil {
		return fmt.Errorf("create notification: %w", err)
	}

	return nil
}

func (r *Repo) fanOutNewPost(ctx context.Context, authorID, postID uuid.UUID) (int64, error) {
	res := r.db.WithContext(ctx).Exec(`
		INSERT INTO notifications (actor_id, recipient_id, type, post_id)
		SELECT @author_id, f.follower_id,  @type, @post_id
		FROM follows f 
		WHERE f.followee_id = @author_id
		ON CONFLICT DO NOTHING
	`,
		sql.Named("author_id", authorID),
		sql.Named("type", string(TypeNewPost)),
		sql.Named("post_id", postID),
	)

	// 貼文在 fan-out 之前就被刪了 → post_id 的 FK 檢查失敗
	if errors.Is(res.Error, gorm.ErrForeignKeyViolated) {
		return 0, errPostGone
	}

	if res.Error != nil {
		return 0, fmt.Errorf("fan out new post (author %s, post %s): %w", authorID, postID, res.Error)
	}

	return res.RowsAffected, nil
}
