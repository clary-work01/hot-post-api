package user

import (
	"context"
	"errors"

	"github.com/google/uuid"
)

var ErrCannotFollowSelf = errors.New("cannot follow yourself")

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
