package auth

import (
	"net/http"

	"github.com/clary-work01/hot-post-api/internal/platform/httpx"
	"github.com/gin-gonic/gin"
)

type loginReq struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// POST /auth/login
func (h *Handler) login(c *gin.Context) {
	var req loginReq
	if err := httpx.DecodeJSON(c, &req); err != nil {
		httpx.InvalidInputError(c, err)
		return
	}

	res, err := h.svc.login(c.Request.Context(), req.Username, req.Password)
	if err != nil {
		writeError(c, err)
		return
	}

	h.setRefreshCookie(c, res.RefreshToken)

	httpx.JSON(c, http.StatusOK, tokenResp{
		AccessToken: res.AccessToken,
	})
}
