package api

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/clary-work01/hot-post-api/internal/domain"
	"github.com/clary-work01/hot-post-api/internal/domain/post"
	"github.com/google/uuid"
)

type PostHandler struct {
	svc *post.Service
}

func NewPostHandler(svc *post.Service) *PostHandler {
	return &PostHandler{
		svc: svc,
	}
}

func (h *PostHandler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /posts", h.create)
	mux.HandleFunc("GET /posts/{id}", h.detail)
}

// ---- request DTO ----
type createPostReq struct {
	Title   string `json:"title"`
	Content string `json:"content"`
}

// ---- response DTO ----
type postResp struct {
	ID        uuid.UUID `json:"id"`
	AuthorID  uuid.UUID `json:"author_id"`
	Title     string    `json:"title"`
	Content   string    `json:"content"`
	LikeCount int       `json:"like_count"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type PostDetailResp struct {
	postResp
	AuthorName string `json:"author_name"`
	IsLiked    bool   `json:"is_liked"`
}

type PostListResp struct {
	List []PostDetailResp `json:"list"`
}

func toPostResp(p *post.Post) postResp {
	return postResp{
		ID:        p.ID,
		AuthorID:  p.AuthorID,
		Title:     p.Title,
		Content:   p.Content,
		LikeCount: p.LikeCount,
		CreatedAt: p.CreatedAt,
		UpdatedAt: p.UpdatedAt,
	}
}

func toPostDetailResp(p *post.PostDetail) PostDetailResp {
	return PostDetailResp{
		postResp:   toPostResp(&p.Post),
		AuthorName: p.AuthorName,
		IsLiked:    p.IsLiked,
	}
}

// POST /posts
func (h *PostHandler) create(w http.ResponseWriter, r *http.Request) {
	var req createPostReq
	if err := decodeJSON(w, r, &req); err != nil {
		writeError(w, r, domain.ErrInvalidInput)
		return
	}

	// ---- Validation ----
	req.Title = strings.TrimSpace(req.Title)
	req.Content = strings.TrimSpace(req.Content)

	if req.Title == "" || req.Content == "" {
		writeError(w, r, domain.ErrInvalidInput)
		return
	}

	// ---- Insert ----
	authorID := uuid.MustParse("d2c5d597-4b6b-4d11-ac6a-06b5b66e1f47") // TODO: 從 JWT 獲得
	p, err := h.svc.Create(
		r.Context(),
		authorID,
		req.Title,
		req.Content,
	)
	if err != nil {
		writeError(w, r, err)
		return
	}

	JSON(w, http.StatusCreated, "post created successfully", toPostResp(p))
}

// GET /posts/{id}
func (h *PostHandler) detail(w http.ResponseWriter, r *http.Request) {
	postID, err := parseID(r)
	if err != nil {
		writeError(w, r, fmt.Errorf("parse id: %w", domain.ErrInvalidInput))
		return
	}

	viewerID := uuid.MustParse("d2c5d597-4b6b-4d11-ac6a-06b5b66e1f47") // TODO: 從 JWT 獲得

	p, err := h.svc.Detail(r.Context(), postID, viewerID)
	if err != nil {
		writeError(w, r, err)
		return
	}

	JSON(w, http.StatusOK, "get post detail success", toPostDetailResp(p))
}
