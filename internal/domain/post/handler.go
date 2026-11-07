package post

import (
	"errors"
	"net/http"
	"time"

	"github.com/clary-work01/hot-post-api/internal/platform/httpx"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type Handler struct {
	svc *Service
}

func NewHandler(svc *Service) *Handler {
	return &Handler{
		svc: svc,
	}
}

func (h *Handler) RegisterRoutes(r gin.IRouter, requireAuth, optionalAuth gin.HandlerFunc) {
	g := r.Group("/posts")

	// 訪客可看；有登入時會帶 is_liked
	pub := g.Group("", optionalAuth)
	pub.GET("/:id", h.detail)
	pub.GET("", h.list)

	// 一定要登入
	priv := g.Group("", requireAuth)
	priv.POST("", h.create)
	priv.PUT("/:id/like", h.like)
	priv.DELETE("/:id/like", h.unlike)

}

func writeError(c *gin.Context, err error) {
	switch {
	case errors.Is(err, errPostNotFound):
		httpx.Error(c, http.StatusNotFound, httpx.ErrorBody{Code: "post_not_found", Message: "貼文不存在"})

	default:
		httpx.ServerError(c, err)
	}
}

// ----  DTO ----
type postResp struct {
	ID        uuid.UUID `json:"id"`
	AuthorID  uuid.UUID `json:"author_id"`
	Title     string    `json:"title"`
	Content   string    `json:"content"`
	LikeCount int       `json:"like_count"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func toPostResp(p *Post) postResp {
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
