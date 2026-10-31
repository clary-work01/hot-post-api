package auth

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"fmt"
	"time"

	"github.com/clary-work01/hot-post-api/internal/domain"
	"github.com/clary-work01/hot-post-api/internal/domain/user"
	"github.com/clary-work01/hot-post-api/internal/token"
	"golang.org/x/crypto/bcrypt"
)

const (
	AccessTTL  = 15 * time.Minute
	RefreshTTL = 7 * 24 * time.Hour
)

type Service struct {
	repo     *Repo
	userRepo *user.Repo
	tokenMgr *token.Manager
}

func NewService(repo *Repo, userRepo *user.Repo, tokenMgr *token.Manager) *Service {
	return &Service{
		repo:     repo,
		userRepo: userRepo,
		tokenMgr: tokenMgr,
	}
}

func (s *Service) Register(ctx context.Context, username, password string) (*user.User, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("hash password: %w", err)
	}

	u := user.User{
		Username: username,
		Password: string(hash),
	}

	if err := s.userRepo.Create(ctx, &u); err != nil {
		return nil, err
	}

	return &u, nil
}

type LoginResult struct {
	AccessToken  string
	RefreshToken string // 原文，只回給前端，DB 只存 hash
}

func (s *Service) Login(ctx context.Context, username, password string) (*LoginResult, error) {
	u, err := s.userRepo.GetUserByUsername(ctx, username)
	if errors.Is(err, domain.ErrNotFound) {
		// 帳號不存在和密碼錯誤回一樣的錯，不讓人猜出哪些帳號存在
		return nil, ErrInvalidCredentials
	}
	if err != nil {
		return nil, err
	}

	err = bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
	if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
		return nil, ErrInvalidCredentials
	}
	if err != nil {
		return nil, fmt.Errorf("compare password: %w", err)
	}

	access, err := s.tokenMgr.Sign(u.ID)
	if err != nil {
		return nil, fmt.Errorf("sign access token: %w", err)
	}

	refresh, err := randomToken()
	if err != nil {
		return nil, fmt.Errorf("generate refresh token: %w", err)
	}

	if err := s.repo.createRefreshToken(ctx, &RefreshToken{
		UserID:    u.ID,
		TokenHash: hashToken(refresh),
		ExpiresAt: time.Now().Add(RefreshTTL),
	}); err != nil {
		return nil, err
	}

	return &LoginResult{
		AccessToken:  access,
		RefreshToken: refresh,
	}, nil
}

// 32 bytes 的 crypto/rand，夠長猜不到
func randomToken() (string, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", fmt.Errorf("read random bytes: %w", err)
	}

	return base64.RawURLEncoding.EncodeToString(b), nil
}

func hashToken(t string) string {
	sum := sha256.Sum256([]byte(t))
	return hex.EncodeToString(sum[:])
}
