package post

import (
	"net/http"

	"github.com/clary-work01/hot-post-api/internal/platform/httpx"
	"github.com/gin-gonic/gin"
)

type postDetailResp struct {
	postResp
	AuthorName string `json:"author_name"`
	IsLiked    bool   `json:"is_liked"`
}

func toPostDetailResp(p *Detail) postDetailResp {
	return postDetailResp{
		postResp:   toPostResp(&p.Post),
		AuthorName: p.AuthorName,
		IsLiked:    p.IsLiked,
	}
}

// GET /posts/{id}
func (h *Handler) detail(c *gin.Context) {
	postID, err := httpx.ParseID(c)
	if err != nil {
		httpx.InvalidInputError(c, err)
		return
	}

	viewerID := httpx.ViewerID(c) // 訪客為 uuid.Nil

	p, err := h.svc.detail(c.Request.Context(), postID, viewerID)
	if err != nil {
		writeError(c, err)
		return
	}

	httpx.JSON(c, http.StatusOK, toPostDetailResp(p))
}
