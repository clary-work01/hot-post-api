package user

import (
	"context"

	"github.com/clary-work01/hot-post-api/internal/domain"
	"github.com/google/uuid"
)

type Service struct {
	repo *Repo
}

func NewService(repo *Repo) *Service {
	return &Service{repo}
}

func (s *Service) Follow(ctx context.Context, followerID, followeeID uuid.UUID) error {
	if followerID == followeeID {
		return ErrCannotFollowSelf
	}

	follow := Follow{
		FollowerID: followerID,
		FolloweeID: followeeID,
	}

	_, err := s.repo.follow(ctx, &follow)
	return err
}

func (s *Service) UnFollow(ctx context.Context, followerID, followeeID uuid.UUID) error {
	if followerID == followeeID {
		return ErrCannotFollowSelf
	}

	_, err := s.repo.unFollow(ctx, followerID, followeeID)
	return err
}

func (s *Service) Profile(ctx context.Context, userID, viewerID uuid.UUID) (*Profile, error) {
	return s.repo.profile(ctx, userID, viewerID)
}

func (s *Service) Followers(ctx context.Context, userID, viewerID uuid.UUID, offset, limit int) ([]FollowItem, error) {
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
			return nil, domain.ErrNotFound
		}
	}
	return items, nil
}

func (s *Service) Followees(ctx context.Context, userID, viewerID uuid.UUID, offset, limit int) ([]FollowItem, error) {
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
			return nil, domain.ErrNotFound
		}
	}
	return items, nil
}
