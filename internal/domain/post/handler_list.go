package post

import (
	"net/http"
	"time"

	"github.com/clary-work01/hot-post-api/internal/platform/httpx"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

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

func toPostSummaryResp(p *Summary) postSummaryResp {
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

func toPostListResp(list []Summary, total int64) postListResp {
	resp := make([]postSummaryResp, len(list))

	for i, summary := range list {
		resp[i] = toPostSummaryResp(&summary)
	}
	return postListResp{
		List:  resp,
		Total: total,
	}
}

// GET /posts?limit=&offset=
func (h *Handler) list(c *gin.Context) {
	page, err := httpx.BindPage(c)
	if err != nil {
		httpx.InvalidInputError(c, err)
		return
	}

	viewerID := httpx.ViewerID(c) // 訪客為 uuid.Nil

	list, total, err := h.svc.list(c.Request.Context(), page.Offset, page.Limit, viewerID)
	if err != nil {
		writeError(c, err)
		return
	}

	httpx.JSON(c, http.StatusOK, toPostListResp(list, total))
}
