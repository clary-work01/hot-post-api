package token

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
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

func (m *Manager) TTL() time.Duration {
	return m.ttl
}
