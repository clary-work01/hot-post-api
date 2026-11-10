package notification

import (
	"context"
	"errors"
	"log/slog"

	"github.com/clary-work01/hot-post-api/internal/platform/workerpool"
	"github.com/google/uuid"
)

type Service struct {
	repo *Repo
	pool *workerpool.Pool
}

func NewService(repo *Repo, pool *workerpool.Pool) *Service {
	return &Service{repo: repo, pool: pool}
}

func (s *Service) NotifyPostLiked(ctx context.Context, actorID, authorID, postID uuid.UUID) error {
	if actorID == authorID {
		return nil // 自己讚自己的貼文，不通知
	}

	n := &Notification{
		ActorID:     actorID,
		RecipientID: authorID,
		PostID:      &postID,
		Type:        TypePostLiked,
	}

	if err := s.repo.create(ctx, n); err != nil {
		return err
	}

	return nil
}

func (s *Service) NotifyFollowed(ctx context.Context, followerID, followeeID uuid.UUID) error {

	n := &Notification{
		ActorID:     followerID,
		RecipientID: followeeID,
		Type:        TypeUserFollowed,
	}

	if err := s.repo.create(ctx, n); err != nil {
		return err
	}

	return nil
}

func (s *Service) NotifyNewPost(ctx context.Context, authorID, postID uuid.UUID) {
	ok := s.pool.Submit(ctx, func(ctx context.Context) error {
		return s.fanOutNewPost(ctx, authorID, postID) // error 由 pool 統一 log
	})

	if !ok {
		slog.WarnContext(ctx, "notification job dropped",
			"post_id", postID,
			"author_id", authorID,
		)
	}
}

func (s *Service) fanOutNewPost(ctx context.Context, authorID, postID uuid.UUID) error {
	_, err := s.repo.fanOutNewPost(ctx, authorID, postID)

	// 貼文已刪除，沒有東西要通知，不算錯誤
	if errors.Is(err, errPostGone) {
		return nil
	}

	return err
}
