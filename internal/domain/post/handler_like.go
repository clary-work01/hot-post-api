package post

import (
	"github.com/clary-work01/hot-post-api/internal/platform/httpx"
	"github.com/gin-gonic/gin"
)

// PUT /posts/{id}/like
func (h *Handler) like(c *gin.Context) {
	postID, err := httpx.ParseID(c)
	if err != nil {
		httpx.InvalidInputError(c, err)
		return
	}

	userID, ok := httpx.CurrentUser(c)
	if !ok {
		return
	}

	if err := h.svc.like(c.Request.Context(), postID, userID); err != nil {
		writeError(c, err)
		return
	}

	httpx.NoContent(c)
}

// DELETE /posts/{id}/like
func (h *Handler) unlike(c *gin.Context) {
	postID, err := httpx.ParseID(c)
	if err != nil {
		httpx.InvalidInputError(c, err)
		return
	}

	userID, ok := httpx.CurrentUser(c)
	if !ok {
		return
	}

	if err := h.svc.unLike(c.Request.Context(), postID, userID); err != nil {
		writeError(c, err)
		return
	}

	httpx.NoContent(c)
}
