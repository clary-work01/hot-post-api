package post

import (
	"net/http"

	"github.com/clary-work01/hot-post-api/internal/platform/httpx"
	"github.com/gin-gonic/gin"
)

type createPostReq struct {
	Title   string `json:"title" binding:"required"`
	Content string `json:"content" binding:"required"`
}

// POST /posts
func (h *Handler) create(c *gin.Context) {
	var req createPostReq
	if err := httpx.DecodeJSON(c, &req); err != nil {
		httpx.InvalidInputError(c, err)
		return
	}

	authorID, ok := httpx.CurrentUser(c)
	if !ok {
		return
	}

	p, err := h.svc.create(
		c.Request.Context(),
		authorID,
		req.Title,
		req.Content,
	)
	if err != nil {
		writeError(c, err)
		return
	}

	httpx.JSON(c, http.StatusCreated, toPostResp(p))
}
