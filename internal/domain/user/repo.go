package user

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

func (r *Repo) Create(ctx context.Context, u *User) error {
	err := r.db.WithContext(ctx).Create(u).Error

	if errors.Is(err, gorm.ErrDuplicatedKey) {
		return ErrUsernameTaken
	}

	if err != nil {
		return fmt.Errorf("create user (%s): %w", u.Username, err)
	}

	return nil
}

func (r *Repo) GetUserByUsername(ctx context.Context, username string) (*User, error) {
	var u User
	err := r.db.WithContext(ctx).Where("username = ?", username).Take(&u).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, domain.ErrNotFound
	}

	if err != nil {
		return nil, fmt.Errorf("get user by username (%s): %w", username, err)
	}

	return &u, nil
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
		return false, fmt.Errorf("follow (%s->%s): %w", follow.FollowerID, follow.FolloweeID, err)
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
		return false, fmt.Errorf("unfollow (%s->%s): %w", followerID, followeeID, err)
	}

	return res.RowsAffected == 1, nil
}

func (r *Repo) profile(ctx context.Context, userID, viewerID uuid.UUID) (*Profile, error) {
	var profile Profile

	res := r.db.WithContext(ctx).
		Table("users AS u").
		Select(
			`u.id, u.username,
			(SELECT COUNT(*) FROM posts p WHERE p.author_id = u.id) AS post_count,
			(SELECT COUNT(*) FROM follows f WHERE f.followee_id = u.id) AS follower_count,
			(SELECT COUNT(*) FROM follows f WHERE f.follower_id = u.id) AS following_count,
			EXISTS (
				SELECT 1 from follows f 
				WHERE follower_id = ? AND followee_id = u.id
			) AS is_following
		`, viewerID).
		Where("u.id = ?", userID).
		Scan(&profile)

	if res.Error != nil {
		return nil, fmt.Errorf("get profile (user %s, viewer %s): %w", userID, viewerID, res.Error)
	}

	// Scan 查不到不會回 ErrRecordNotFound，要自己看筆數
	if res.RowsAffected == 0 {
		return nil, domain.ErrNotFound
	}

	return &profile, nil
}

func (r *Repo) followers(ctx context.Context, followeeID, viewerID uuid.UUID, offset, limit int) ([]FollowItem, error) {
	var items []FollowItem

	err := r.db.WithContext(ctx).
		Raw(`
			SELECT u.id, u.username, f.created_at AS followed_at,
			EXISTS(
				SELECT 1 FROM follows mf
				WHERE mf.follower_id = ? AND mf.followee_id = u.id
			) AS is_following
			FROM users u
			JOIN follows f ON u.id = f.follower_id
			WHERE f.followee_id = ?
			ORDER BY f.created_at DESC, f.follower_id DESC
			LIMIT ?
			OFFSET ?
		`, viewerID, followeeID, limit, offset).Scan(&items).Error
	if err != nil {
		return nil, fmt.Errorf("list followers (user %s): %w", followeeID, err)
	}

	return items, err
}

func (r *Repo) followees(ctx context.Context, followerID, viewerID uuid.UUID, offset, limit int) ([]FollowItem, error) {
	var items []FollowItem

	err := r.db.WithContext(ctx).
		Raw(`
			SELECT u.id, u.username, f.created_at AS followed_at,
			EXISTS(
				SELECT 1 FROM follows mf
				WHERE mf.follower_id = ? AND mf.followee_id = u.id
			) AS is_following
			FROM users u
			JOIN follows f ON u.id = f.followee_id
			WHERE f.follower_id = ?
			ORDER BY f.created_at DESC, f.followee_id DESC
			LIMIT ?
			OFFSET ?
		`, viewerID, followerID, limit, offset).Scan(&items).Error

	if err != nil {
		return nil, fmt.Errorf("list followees (user %s): %w", followerID, err)
	}

	return items, err
}

func (r *Repo) userExists(ctx context.Context, id uuid.UUID) (bool, error) {
	var exists bool
	err := r.db.WithContext(ctx).
		Raw(
			`SELECT EXISTS (
				SELECT 1 FROM users
				WHERE id = ?
			)`, id).
		Scan(&exists).Error

	if err != nil {
		return false, fmt.Errorf("check user exists (%s): %w", id, err)
	}

	return exists, nil
}
