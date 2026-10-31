package auth

import (
	"context"
	"fmt"

	"github.com/clary-work01/hot-post-api/internal/domain/user"
	"golang.org/x/crypto/bcrypt"
)

type Service struct {
	userRepo *user.Repo
}

func NewService(userRepo *user.Repo) *Service {
	return &Service{userRepo}
}

func (s *Service) Register(ctx context.Context, username, password string) (*user.User, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("register (username %s): generate from password: %w", username, err)
	}

	u := user.User{
		Username: username,
		Password: string(hash),
	}

	if err := s.userRepo.Create(ctx, &u); err != nil {
		return nil, fmt.Errorf("register: %w", err)
	}

	return &u, nil
}
