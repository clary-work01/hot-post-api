package api

import (
	"net/http"
	"time"

	"github.com/clary-work01/hot-post-api/internal/domain/user"
	"github.com/google/uuid"
)

type UserHandler struct {
	svc         *user.Service
	requireAuth Middleware
}

func NewUserHandler(svc *user.Service, requireAuth Middleware) *UserHandler {
	return &UserHandler{svc: svc, requireAuth: requireAuth}
}

func (h *UserHandler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("PUT /users/{id}/follow", h.requireAuth(h.follow))
	mux.HandleFunc("DELETE /users/{id}/follow", h.requireAuth(h.unfollow))
	mux.HandleFunc("GET /users/{id}", h.requireAuth(h.profile))
	mux.HandleFunc("GET /users/{id}/followers", h.requireAuth(h.followers))
	mux.HandleFunc("GET /users/{id}/followees", h.requireAuth(h.followees))
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

// PUT /users/{id}/follow
func (h *UserHandler) follow(w http.ResponseWriter, r *http.Request) {
	followeeID, err := parseID(r)
	if err != nil {
		errorJSON(w, http.StatusBadRequest, codeInvalidInput, "參數id錯誤")
		return
	}

	followerID, ok := currentUser(w, r)
	if !ok {
		return
	}

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
		errorJSON(w, http.StatusBadRequest, codeInvalidInput, "參數id錯誤")
		return
	}

	followerID, ok := currentUser(w, r)
	if !ok {
		return
	}

	err = h.svc.UnFollow(r.Context(), followerID, followeeID)
	if err != nil {
		writeError(w, r, err)
		return
	}

	JSON(w, http.StatusOK, "unfollow user success", nil)
}

// GET /users/{id}
func (h *UserHandler) profile(w http.ResponseWriter, r *http.Request) {
	userID, err := parseID(r)
	if err != nil {
		errorJSON(w, http.StatusBadRequest, codeInvalidInput, "參數id錯誤")
		return
	}
	viewerID, ok := currentUser(w, r)
	if !ok {
		return
	}

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
		errorJSON(w, http.StatusBadRequest, codeInvalidInput, "參數id錯誤")
		return
	}

	viewerID, ok := currentUser(w, r)
	if !ok {
		return
	}

	q := r.URL.Query()
	offset, ok := parseOffset(q.Get("offset"))
	if !ok {
		errorJSON(w, http.StatusBadRequest, codeInvalidInput, "參數offset錯誤")
		return
	}

	limit, ok := parseLimit(q.Get("limit"))
	if !ok {
		errorJSON(w, http.StatusBadRequest, codeInvalidInput, "參數limit錯誤")
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
		errorJSON(w, http.StatusBadRequest, codeInvalidInput, "參數id錯誤")
		return
	}
	viewerID, ok := currentUser(w, r)
	if !ok {
		return
	}

	q := r.URL.Query()
	offset, ok := parseOffset(q.Get("offset"))
	if !ok {
		errorJSON(w, http.StatusBadRequest, codeInvalidInput, "參數offset錯誤")
		return
	}

	limit, ok := parseLimit(q.Get("limit"))
	if !ok {
		errorJSON(w, http.StatusBadRequest, codeInvalidInput, "參數limit錯誤")
		return
	}

	items, err := h.svc.Followees(r.Context(), userID, viewerID, offset, limit)
	if err != nil {
		writeError(w, r, err)
		return
	}

	JSON(w, http.StatusOK, "list followees success", toFollowListResp(items))
}
