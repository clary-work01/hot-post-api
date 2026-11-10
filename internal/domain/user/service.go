package user

import (
	"context"
	"log/slog"

	"github.com/google/uuid"
)

type Service struct {
	repo     *Repo
	notifier followNotifier
}

type followNotifier interface {
	NotifyFollowed(ctx context.Context, followerID, followeeId uuid.UUID) error
}

func NewService(repo *Repo, notifier followNotifier) *Service {
	return &Service{repo: repo, notifier: notifier}
}

func (s *Service) follow(ctx context.Context, followerID, followeeID uuid.UUID) error {
	if followerID == followeeID {
		return errCannotFollowSelf
	}

	follow := Follow{
		FollowerID: followerID,
		FolloweeID: followeeID,
	}

	inserted, err := s.repo.follow(ctx, &follow)
	if err != nil {
		return err
	}

	if !inserted {
		return nil
	}

	if err := s.notifier.NotifyFollowed(ctx, followerID, followeeID); err != nil {
		// 通知只是附帶效果，失敗了也不能讓follow失敗。
		// 錯誤在這裡被吞掉、不再往上傳，所以這裡就是它的邊界，在這裡 log 一次
		slog.WarnContext(ctx, "notify post liked failed",
			"error", err,
			"follower_id", followerID,
			"followee_id", followeeID,
		)
	}

	return nil
}

func (s *Service) unFollow(ctx context.Context, followerID, followeeID uuid.UUID) error {
	if followerID == followeeID {
		return errCannotFollowSelf
	}

	_, err := s.repo.unFollow(ctx, followerID, followeeID)
	return err
}

func (s *Service) profile(ctx context.Context, userID, viewerID uuid.UUID) (*Profile, error) {
	return s.repo.profile(ctx, userID, viewerID)
}

func (s *Service) followers(ctx context.Context, userID, viewerID uuid.UUID, offset, limit int) ([]FollowItem, error) {
	// 先樂觀查，空了才確認
	items, err := s.repo.followers(ctx, userID, viewerID, offset, limit)
	if err != nil {
		return nil, err
	}

	// 只有結果為空陣列時才需要分辨「沒粉絲」和「使用者不存在」
	if len(items) == 0 {
		exists, err := s.repo.userExists(ctx, userID)
		if err != nil {
			return nil, err
		}
		if !exists {
			return nil, ErrUserNotFound
		}
	}
	return items, nil
}

func (s *Service) followees(ctx context.Context, userID, viewerID uuid.UUID, offset, limit int) ([]FollowItem, error) {
	items, err := s.repo.followees(ctx, userID, viewerID, offset, limit)
	if err != nil {
		return nil, err
	}
	// 只有結果為空陣列時才需要分辨「沒關注」和「使用者不存在」
	if len(items) == 0 {
		exists, err := s.repo.userExists(ctx, userID)
		if err != nil {
			return nil, err
		}
		if !exists {
			return nil, ErrUserNotFound
		}
	}
	return items, nil
}
