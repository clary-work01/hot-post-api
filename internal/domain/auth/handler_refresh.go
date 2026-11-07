package auth

import (
	"errors"
	"net/http"

	"github.com/clary-work01/hot-post-api/internal/platform/httpx"
	"github.com/gin-gonic/gin"
)

// POST /auth/refresh
func (h *Handler) refresh(c *gin.Context) {
	cookie, err := c.Request.Cookie(refreshCookieName)
	if err != nil {
		writeError(c, errInvalidRefreshToken)
		return
	}

	res, err := h.svc.refresh(c.Request.Context(), cookie.Value)
	if errors.Is(err, errInvalidRefreshToken) {
		h.clearRefreshCookie(c)
	}
	if err != nil {
		writeError(c, err)
		return
	}

	h.setRefreshCookie(c, res.RefreshToken)

	httpx.JSON(c, http.StatusOK, tokenResp{
		AccessToken: res.AccessToken,
	})
}
