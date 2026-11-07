package auth

import (
	"errors"
	"net/http"

	"github.com/clary-work01/hot-post-api/internal/domain/user"
	"github.com/clary-work01/hot-post-api/internal/platform/httpx"
	"github.com/gin-gonic/gin"
)

type Handler struct {
	svc            *Service
	isCookieSecure bool
}

func NewHandler(svc *Service, isCookieSecure bool) *Handler {
	return &Handler{
		svc:            svc,
		isCookieSecure: isCookieSecure,
	}
}

func (h *Handler) RegisterRoutes(r gin.IRouter) {
	g := r.Group("/auth")

	g.POST("/register", h.register)
	g.POST("/login", h.login)
	g.POST("/refresh", h.refresh)
	g.POST("/logout", h.logout)
}

// ---- 共用  ----

const refreshCookieName = "refresh_token"

type tokenResp struct {
	AccessToken string `json:"access_token"`
}

func writeError(c *gin.Context, err error) {
	switch {
	case errors.Is(err, errInvalidRefreshToken):
		httpx.Error(c, http.StatusUnauthorized, httpx.ErrorBody{Code: "invalid_refresh_token", Message: "請重新登入"})

	case errors.Is(err, errInvalidCredentials):
		httpx.Error(c, http.StatusUnauthorized, httpx.ErrorBody{Code: "invalid_credentials", Message: "帳密錯誤"})

	// 來自 user feature 的錯誤
	case errors.Is(err, user.ErrUsernameTaken):
		httpx.Error(c, http.StatusConflict, httpx.ErrorBody{Code: "username_taken", Message: "用戶名已存在"})

	case errors.Is(err, user.ErrUserNotFound):
		httpx.Error(c, http.StatusConflict, httpx.ErrorBody{Code: "user_not_found", Message: "用戶不存在"})

	default:
		httpx.ServerError(c, err)
	}
}

func (h *Handler) setRefreshCookie(c *gin.Context, refreshToken string) {
	http.SetCookie(c.Writer, &http.Cookie{
		Name:  refreshCookieName,
		Value: refreshToken,
		// 只有 /auth/... 開頭的請求會帶上
		Path: "/auth",
		// 存活秒數，7 天後瀏覽器自動刪除
		MaxAge: int(RefreshTTL.Seconds()),
		// JS 讀不到
		HttpOnly: true,
		// 只在 HTTPS 連線時才送出
		Secure: h.isCookieSecure,
		// 只有從你自己網站發出的請求才帶
		SameSite: http.SameSiteStrictMode,
	})
}

func (h *Handler) clearRefreshCookie(c *gin.Context) {
	http.SetCookie(c.Writer, &http.Cookie{
		Name:     refreshCookieName,
		Value:    "",
		Path:     "/auth",
		MaxAge:   -1,
		HttpOnly: true,
		Secure:   h.isCookieSecure,
		SameSite: http.SameSiteStrictMode,
	})
}
