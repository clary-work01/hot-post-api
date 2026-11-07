package user

import (
	"net/http"

	"github.com/clary-work01/hot-post-api/internal/platform/httpx"
	"github.com/gin-gonic/gin"
)

// GET /users/:id/followers?offset=&limit=
func (h *Handler) followers(c *gin.Context) {
	userID, err := httpx.ParseID(c)
	if err != nil {
		httpx.InvalidInputError(c, err)
		return
	}

	page, err := httpx.BindPage(c)
	if err != nil {
		httpx.InvalidInputError(c, err)
		return
	}

	viewerID := httpx.ViewerID(c) // 訪客為 uuid.Nil

	items, err := h.svc.followers(c.Request.Context(), userID, viewerID, page.Offset, page.Limit)
	if err != nil {
		writeError(c, err)
		return
	}

	httpx.JSON(c, http.StatusOK, toFollowListResp(items))
}
