package api

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/google/uuid"
)

const (
	defaultLimit = 20
	maxLimit     = 100
)

func parseID(r *http.Request) (uuid.UUID, error) {
	id, err := uuid.Parse(r.PathValue("id"))

	if err != nil {
		return uuid.Nil, err
	}
	// id == uuid.Nil: 00000000-0000-0000-0000-000000000000
	if id == uuid.Nil {
		return uuid.Nil, errors.New("nil uuid")
	}
	return id, nil
}

func parseOffset(s string) (offset int, ok bool) {
	if s == "" {
		return 0, true
	}

	n, err := strconv.Atoi(s)
	if err != nil || n < 0 {
		return 0, false
	}
	return n, true
}

func parseLimit(s string) (limit int, ok bool) {
	if s == "" {
		return defaultLimit, true
	}

	n, err := strconv.Atoi(s)
	if err != nil || n < 0 {
		return 0, false
	}
	if n > maxLimit {
		n = maxLimit
	}
	return n, true
}
