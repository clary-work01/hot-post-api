package post

import (
	"context"

	"github.com/google/uuid"
)

type Service struct {
	repo *Repo
}

func NewService(repo *Repo) *Service {
	return &Service{
		repo: repo,
	}
}

func (s *Service) Create(ctx context.Context, authorID uuid.UUID, title string, content string) (*Post, error) {
	p := &Post{
		AuthorID: authorID,
		Title:    title,
		Content:  content,
	}

	if err := s.repo.create(ctx, p); err != nil {
		return nil, err
	}

	return p, nil
}

func (s *Service) Detail(ctx context.Context, postID, viewerID uuid.UUID) (*Detail, error) {
	return s.repo.detail(ctx, postID, viewerID)
}

func (s *Service) List(ctx context.Context, offset, limit int, viewerID uuid.UUID) (list []Summary, total int64, err error) {
	return s.repo.list(ctx, offset, limit, viewerID)
}

func (s *Service) Like(ctx context.Context, postID, userID uuid.UUID) error {
	postLike := PostLike{
		PostID: postID,
		UserID: userID,
	}

	_, err := s.repo.like(ctx, &postLike)
	if err != nil {
		return err
	}

	return nil
}

func (s *Service) UnLike(ctx context.Context, postID, userID uuid.UUID) error {
	_, err := s.repo.unlike(ctx, postID, userID)
	if err != nil {
		return err
	}

	return nil
}
