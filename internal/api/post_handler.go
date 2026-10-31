package api

import (
	"errors"
	"net/http"
	"strings"
	"time"

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
	mux.HandleFunc("GET /posts", h.list)

}

// ---- request DTO ----
type createPostReq struct {
	Title   string `json:"title"`
	Content string `json:"content"`
}

func (req *createPostReq) validate() error {
	req.Title = strings.TrimSpace(req.Title)
	req.Content = strings.TrimSpace(req.Content)

	if req.Title == "" || req.Content == "" {
		return errors.New("參數title、content不得為空字串")
	}
	return nil
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

type postDetailResp struct {
	postResp
	AuthorName string `json:"author_name"`
	IsLiked    bool   `json:"is_liked"`
}

type postSummaryResp struct {
	ID         uuid.UUID `json:"id"`
	AuthorID   uuid.UUID `json:"author_id"`
	AuthorName string    `json:"author_name"`
	Title      string    `json:"title"`
	Excerpt    string    `json:"excerpt"`
	LikeCount  int       `json:"like_count"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
	IsLiked    bool      `json:"is_liked"`
}

type postListResp struct {
	List  []postSummaryResp `json:"list"`
	Total int64             `json:"total"`
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

func toPostDetailResp(p *post.Detail) postDetailResp {
	return postDetailResp{
		postResp:   toPostResp(&p.Post),
		AuthorName: p.AuthorName,
		IsLiked:    p.IsLiked,
	}
}

func toPostSummaryResp(p *post.Summary) postSummaryResp {
	return postSummaryResp{
		ID:         p.ID,
		AuthorID:   p.AuthorID,
		AuthorName: p.AuthorName,
		Title:      p.Title,
		Excerpt:    p.Excerpt,
		LikeCount:  p.LikeCount,
		CreatedAt:  p.CreatedAt,
		UpdatedAt:  p.UpdatedAt,
		IsLiked:    p.IsLiked,
	}
}

func toPostListResp(list []post.Summary, total int64) postListResp {
	resp := make([]postSummaryResp, len(list))

	for i, summary := range list {
		resp[i] = toPostSummaryResp(&summary)
	}
	return postListResp{
		List:  resp,
		Total: total,
	}
}

// POST /posts
func (h *PostHandler) create(w http.ResponseWriter, r *http.Request) {
	var req createPostReq
	if err := decodeJSON(w, r, &req); err != nil {
		JSON(w, http.StatusBadRequest, "參數錯誤", nil)
		return
	}

	if err := req.validate(); err != nil {
		JSON(w, http.StatusBadRequest, err.Error(), nil)
		return
	}

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
		JSON(w, http.StatusBadRequest, "參數id錯誤", nil)
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

// GET /posts?limit=&offset=
func (h *PostHandler) list(w http.ResponseWriter, r *http.Request) {
	viewerID := uuid.MustParse("d2c5d597-4b6b-4d11-ac6a-06b5b66e1f47") // TODO: 從 JWT 獲得

	q := r.URL.Query()
	offset, ok := parseOffset(q.Get("offset"))
	if !ok {
		JSON(w, http.StatusBadRequest, "參數offset錯誤", nil)
		return
	}

	limit, ok := parseLimit(q.Get("limit"))
	if !ok {
		JSON(w, http.StatusBadRequest, "參數limit錯誤", nil)
		return
	}

	list, total, err := h.svc.List(r.Context(), offset, limit, viewerID)
	if err != nil {
		writeError(w, r, err)
		return
	}

	JSON(w, http.StatusOK, "get post list success", toPostListResp(list, total))
}
