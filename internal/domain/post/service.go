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

func (s *Service) Detail(ctx context.Context, postID, viewerID uuid.UUID) (*PostDetail, error) {
	return s.repo.detail(ctx, postID, viewerID)
}
