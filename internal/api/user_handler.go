package api

import (
	"fmt"
	"net/http"
	"time"

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
	mux.HandleFunc("GET /users/{id}", h.profile)
	mux.HandleFunc("GET /users/{id}/followers", h.followers)
	mux.HandleFunc("GET /users/{id}/followees", h.followees)
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

// ---- DTO ----
type ProfileResp struct {
	ID             uuid.UUID `json:"id"`
	Username       string    `json:"username"`
	PostCount      int       `json:"post_count"`
	FollowerCount  int       `json:"follower_count"`
	FollowingCount int       `json:"following_count"`
	IsFollowing    bool      `json:"is_following"`
}

func toProfileResp(p *user.Profile) *ProfileResp {
	return &ProfileResp{
		ID:             p.ID,
		Username:       p.Username,
		PostCount:      p.PostCount,
		FollowerCount:  p.FollowerCount,
		FollowingCount: p.FollowingCount,
		IsFollowing:    p.IsFollowing,
	}
}

type FollowItemResp struct {
	ID          uuid.UUID `json:"id"`
	Username    string    `json:"username"`
	FollowedAt  time.Time `json:"followed_at"`
	IsFollowing bool      `json:"is_following"`
}

type FollowListResp struct {
	List []FollowItemResp `json:"list"`
}

func toFollowListResp(items []user.FollowItem) *FollowListResp {
	followItems := make([]FollowItemResp, len(items))
	for i, item := range items {
		f := FollowItemResp{
			ID:          item.ID,
			Username:    item.Username,
			FollowedAt:  item.FollowedAt,
			IsFollowing: item.IsFollowing,
		}
		followItems[i] = f
	}

	return &FollowListResp{List: followItems}
}

// GET /users/{id}
func (h *UserHandler) profile(w http.ResponseWriter, r *http.Request) {
	userID, err := parseID(r)
	if err != nil {
		writeError(w, r, fmt.Errorf("parse id %w", domain.ErrInvalidInput))
		return
	}
	viewerID := uuid.MustParse("d2c5d597-4b6b-4d11-ac6a-06b5b66e1f47") // TODO: 從 JWT 獲得

	profile, err := h.svc.Profile(r.Context(), userID, viewerID)
	if err != nil {
		writeError(w, r, err)
		return
	}

	JSON(w, http.StatusOK, "get profile success", toProfileResp(profile))
}

// GET /users/{id}/followers?offset=&limit=
func (h *UserHandler) followers(w http.ResponseWriter, r *http.Request) {
	userID, err := parseID(r)
	if err != nil {
		writeError(w, r, fmt.Errorf("parse id %w", domain.ErrInvalidInput))
		return
	}
	viewerID := uuid.MustParse("d2c5d597-4b6b-4d11-ac6a-06b5b66e1f47") // TODO: 從 JWT 獲得

	q := r.URL.Query()
	offset, ok := parseOffset(q.Get("offset"))
	if !ok {
		writeError(w, r, fmt.Errorf("parse offset %w", domain.ErrInvalidInput))
		return
	}

	limit, ok := parseLimit(q.Get("limit"))
	if !ok {
		writeError(w, r, fmt.Errorf("parse limit %w", domain.ErrInvalidInput))
		return
	}

	items, err := h.svc.Followers(r.Context(), userID, viewerID, offset, limit)
	if err != nil {
		writeError(w, r, err)
		return
	}

	JSON(w, http.StatusOK, "list followers success", toFollowListResp(items))
}

// GET /users/{id}/followees?offset=&limit=
func (h *UserHandler) followees(w http.ResponseWriter, r *http.Request) {
	userID, err := parseID(r)
	if err != nil {
		writeError(w, r, fmt.Errorf("parse id %w", domain.ErrInvalidInput))
		return
	}
	viewerID := uuid.MustParse("d2c5d597-4b6b-4d11-ac6a-06b5b66e1f47") // TODO: 從 JWT 獲得

	q := r.URL.Query()
	offset, ok := parseOffset(q.Get("offset"))
	if !ok {
		writeError(w, r, fmt.Errorf("parse offset %w", domain.ErrInvalidInput))
		return
	}

	limit, ok := parseLimit(q.Get("limit"))
	if !ok {
		writeError(w, r, fmt.Errorf("parse limit %w", domain.ErrInvalidInput))
		return
	}

	items, err := h.svc.Followees(r.Context(), userID, viewerID, offset, limit)
	if err != nil {
		writeError(w, r, err)
		return
	}

	JSON(w, http.StatusOK, "list followees success", toFollowListResp(items))
}
