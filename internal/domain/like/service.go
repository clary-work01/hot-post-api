package like

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
