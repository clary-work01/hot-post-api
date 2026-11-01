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

	refresh, refreshHash := newRefreshToken()

	if err := s.repo.createRefreshToken(ctx, &RefreshToken{
		UserID:    u.ID,
		TokenHash: refreshHash,
		ExpiresAt: time.Now().Add(RefreshTTL),
	}); err != nil {
		return nil, err
	}

	return &LoginResult{
		AccessToken:  access,
		RefreshToken: refresh,
	}, nil
}

type RefreshResult struct {
	AccessToken  string
	RefreshToken string
}

func (s *Service) Refresh(ctx context.Context, oldRefresh string) (*RefreshResult, error) {
	if oldRefresh == "" {
		return nil, ErrInvalidRefreshToken
	}

	newRefresh, newRefreshHash := newRefreshToken()

	userID, err := s.repo.rotateRefreshToken(ctx, hashToken(oldRefresh), &RefreshToken{
		TokenHash: newRefreshHash,
		ExpiresAt: time.Now().Add(RefreshTTL),
	})
	if err != nil {
		return nil, err
	}

	newAccess, err := s.tokenMgr.Sign(userID)
	if err != nil {
		return nil, fmt.Errorf("sign access token: %w", err)
	}

	return &RefreshResult{
		AccessToken:  newAccess,
		RefreshToken: newRefresh,
	}, nil
}

func (s *Service) Logout(ctx context.Context, refresh string) error {
	if refresh == "" {
		return ErrInvalidRefreshToken
	}

	return s.repo.revokeRefreshToken(ctx, hashToken(refresh))
}

func newRefreshToken() (refresh, refreshHash string) {
	refresh = generateToken()
	refreshHash = hashToken(refresh)
	return refresh, refreshHash
}

// 32 bytes 的 crypto/rand，夠長猜不到
func generateToken() string {
	b := make([]byte, 32)
	// Go 1.24 以後 crypto/rand.Read 保證不回傳錯誤（失敗會直接讓程式終止）
	rand.Read(b)
	token := base64.RawURLEncoding.EncodeToString(b)

	return token
}

func hashToken(t string) string {
	sum := sha256.Sum256([]byte(t))
	return hex.EncodeToString(sum[:])
}
