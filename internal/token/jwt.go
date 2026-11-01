package token

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

var (
	ErrInvalidToken = errors.New("invalid token")
	ErrExpiredToken = errors.New("token expired")
)

type Manager struct {
	secret []byte
	ttl    time.Duration
}

func NewManager(secret []byte, ttl time.Duration) *Manager {
	return &Manager{secret, ttl}
}

func (m *Manager) Sign(userID uuid.UUID) (string, error) {
	now := time.Now()

	claims := jwt.RegisteredClaims{
		Subject:   userID.String(),
		IssuedAt:  jwt.NewNumericDate(now),
		ExpiresAt: jwt.NewNumericDate(now.Add(m.ttl)),
	}

	return jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString(m.secret)
}

func (m *Manager) Verify(raw string) (uuid.UUID, error) {
	claims := &jwt.RegisteredClaims{}
	_, err := jwt.ParseWithClaims(raw, claims,
		func(t *jwt.Token) (any, error) {
			return m.secret, nil
		},
		jwt.WithValidMethods([]string{"HS256"}),
		jwt.WithExpirationRequired(),
	)
	if errors.Is(err, jwt.ErrTokenExpired) {
		return uuid.Nil, ErrExpiredToken
	}
	if err != nil {
		return uuid.Nil, ErrInvalidToken
	}

	id, err := uuid.Parse(claims.Subject)
	if err != nil {
		return uuid.Nil, ErrInvalidToken
	}
	return id, nil
}

func (m *Manager) TTL() time.Duration {
	return m.ttl
}
