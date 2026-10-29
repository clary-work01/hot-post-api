package api

import (
	"fmt"
	"net/http"

	"github.com/clary-work01/hot-post-api/internal/domain"
	"github.com/clary-work01/hot-post-api/internal/domain/user"
	"github.com/google/uuid"
)

type UserHandler struct {
	svc *user.Service
}

func NewUserHandler(svc *user.Service) *UserHandler {
	return &UserHandler{svc}
}

func (h *UserHandler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("PUT /users/{id}/follow", h.follow)
	mux.HandleFunc("DELETE /users/{id}/follow", h.unfollow)
}

// PUT /users/{id}/follow
func (h *UserHandler) follow(w http.ResponseWriter, r *http.Request) {
	followeeID, err := parseID(r)
	if err != nil {
		writeError(w, r, fmt.Errorf("parse id %w", domain.ErrInvalidInput))
		return
	}

	followerID := uuid.MustParse("d2c5d597-4b6b-4d11-ac6a-06b5b66e1f47") // TODO: 從 JWT 獲得

	err = h.svc.Follow(r.Context(), followerID, followeeID)
	if err != nil {
		writeError(w, r, err)
		return
	}

	JSON(w, http.StatusOK, "follow user success", nil)
}

// DELETE /users/{id}/unfollow
func (h *UserHandler) unfollow(w http.ResponseWriter, r *http.Request) {
	followeeID, err := parseID(r)
	if err != nil {
		writeError(w, r, fmt.Errorf("parse id %w", domain.ErrInvalidInput))
		return
	}

	followerID := uuid.MustParse("d2c5d597-4b6b-4d11-ac6a-06b5b66e1f47") // TODO: 從 JWT 獲得

	err = h.svc.UnFollow(r.Context(), followerID, followeeID)
	if err != nil {
		writeError(w, r, err)
		return
	}

	JSON(w, http.StatusOK, "unfollow user success", nil)
}
