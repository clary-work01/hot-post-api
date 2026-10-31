package api

import (
	"net/http"

	"github.com/clary-work01/hot-post-api/internal/domain/like"
	"github.com/google/uuid"
)

type LikeHandler struct {
	svc *like.Service
}

func NewLikeHandler(svc *like.Service) *LikeHandler {
	return &LikeHandler{svc}
}

func (h *LikeHandler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("PUT /posts/{id}/like", h.like)
	mux.HandleFunc("DELETE /posts/{id}/like", h.unlike)
}

func (h *LikeHandler) like(w http.ResponseWriter, r *http.Request) {
	postID, err := parseID(r)
	if err != nil {
		JSON(w, http.StatusBadRequest, "參數id錯誤", nil)
		return
	}

	userID := uuid.MustParse("d2c5d597-4b6b-4d11-ac6a-06b5b66e1f47") // TODO: 從 JWT 獲得

	if err := h.svc.Like(r.Context(), postID, userID); err != nil {
		writeError(w, r, err)
		return
	}

	JSON(w, http.StatusOK, "like post success", nil)
}

func (h *LikeHandler) unlike(w http.ResponseWriter, r *http.Request) {
	postID, err := parseID(r)
	if err != nil {
		JSON(w, http.StatusBadRequest, "參數id錯誤", nil)
		return
	}

	userID := uuid.MustParse("d2c5d597-4b6b-4d11-ac6a-06b5b66e1f47") // TODO: 從 JWT 獲得

	if err := h.svc.UnLike(r.Context(), postID, userID); err != nil {
		writeError(w, r, err)
		return
	}

	JSON(w, http.StatusOK, "unlike post success", nil)
}
