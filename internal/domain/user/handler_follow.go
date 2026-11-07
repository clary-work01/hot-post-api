package user

import (
	"github.com/clary-work01/hot-post-api/internal/platform/httpx"
	"github.com/gin-gonic/gin"
)

// PUT /users/:id/follow
func (h *Handler) follow(c *gin.Context) {
	followeeID, err := httpx.ParseID(c)
	if err != nil {
		httpx.InvalidInputError(c, err)
		return
	}

	followerID, ok := httpx.CurrentUser(c)
	if !ok {
		return
	}

	err = h.svc.follow(c.Request.Context(), followerID, followeeID)
	if err != nil {
		writeError(c, err)
		return
	}

	httpx.NoContent(c)
}

// DELETE /users/:id/follow
func (h *Handler) unfollow(c *gin.Context) {
	followeeID, err := httpx.ParseID(c)
	if err != nil {
		httpx.InvalidInputError(c, err)
		return
	}

	followerID, ok := httpx.CurrentUser(c)
	if !ok {
		return
	}

	err = h.svc.unFollow(c.Request.Context(), followerID, followeeID)
	if err != nil {
		writeError(c, err)
		return
	}

	httpx.NoContent(c)
}
