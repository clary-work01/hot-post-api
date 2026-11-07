package user

import (
	"context"

	"github.com/google/uuid"
)

type Service struct {
	repo *Repo
}

func NewService(repo *Repo) *Service {
	return &Service{repo}
}

func (s *Service) follow(ctx context.Context, followerID, followeeID uuid.UUID) error {
	if followerID == followeeID {
		return errCannotFollowSelf
	}

	follow := Follow{
		FollowerID: followerID,
		FolloweeID: followeeID,
	}

	_, err := s.repo.follow(ctx, &follow)
	return err
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
