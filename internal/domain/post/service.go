package post

import (
	"context"
	"log/slog"

	"github.com/google/uuid"
)

type Service struct {
	repo     *Repo
	notifier Notifier
}

type Notifier interface {
	NotifyPostLiked(ctx context.Context, actorID, authorID, postID uuid.UUID) error
	NotifyNewPost(ctx context.Context, authorID, postID uuid.UUID) // 非同步，只負責排入佇列，不回傳 error
}

func NewService(repo *Repo, notifier Notifier) *Service {
	return &Service{
		repo:     repo,
		notifier: notifier,
	}
}

func (s *Service) create(ctx context.Context, authorID uuid.UUID, title string, content string) (*Post, error) {
	p := &Post{
		AuthorID: authorID,
		Title:    title,
		Content:  content,
	}

	postID, err := s.repo.create(ctx, p)
	if err != nil {
		return nil, err
	}
	slog.DebugContext(ctx, "1", "postID", postID)

	s.notifier.NotifyNewPost(ctx, authorID, postID)

	return p, nil
}

func (s *Service) detail(ctx context.Context, postID, viewerID uuid.UUID) (*Detail, error) {
	return s.repo.detail(ctx, postID, viewerID)
}

func (s *Service) list(ctx context.Context, offset, limit int, viewerID uuid.UUID) (list []Summary, total int64, err error) {
	return s.repo.list(ctx, offset, limit, viewerID)
}

func (s *Service) like(ctx context.Context, postID, userID uuid.UUID) error {
	postLike := PostLike{
		PostID: postID,
		UserID: userID,
	}

	authorID, inserted, err := s.repo.like(ctx, &postLike)

	if err != nil {
		return err
	}

	if !inserted {
		return nil
	}

	err = s.notifier.NotifyPostLiked(ctx, userID, authorID, postID)
	// 通知只是附帶效果，失敗了也不能讓按讚失敗。
	// 錯誤在這裡被吞掉、不再往上傳，所以這裡就是它的邊界，在這裡 log 一次
	if err != nil {
		slog.WarnContext(ctx, "notify post liked failed",
			"error", err,
			"actor_id", userID,
			"author_id", authorID,
			"post_id", postID,
		)
	}

	return nil
}

func (s *Service) unLike(ctx context.Context, postID, userID uuid.UUID) error {
	_, err := s.repo.unlike(ctx, postID, userID)
	if err != nil {
		return err
	}

	return nil
}
