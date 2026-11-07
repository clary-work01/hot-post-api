package auth

import (
	"github.com/clary-work01/hot-post-api/internal/platform/httpx"
	"github.com/gin-gonic/gin"
)

func (h *Handler) logout(c *gin.Context) {
	// 不管後面成不成功，瀏覽器端的 cookie 都先清掉
	h.clearRefreshCookie(c)

	cookie, err := c.Request.Cookie(refreshCookieName)
	if err != nil {
		// 沒帶 cookie：本來就沒登入，直接當成功
		httpx.NoContent(c)
		return
	}

	if err := h.svc.logout(c.Request.Context(), cookie.Value); err != nil {
		writeError(c, err) // 只有 DB 錯誤會到這 → 500
		return
	}

	httpx.NoContent(c)
}
