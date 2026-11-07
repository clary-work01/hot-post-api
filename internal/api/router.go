package api

import (
	"log/slog"
	"net/http"

	"github.com/clary-work01/hot-post-api/internal/api/middleware"
	"github.com/clary-work01/hot-post-api/internal/domain/auth"
	"github.com/clary-work01/hot-post-api/internal/domain/post"
	"github.com/clary-work01/hot-post-api/internal/domain/user"
	"github.com/gin-gonic/gin"
)

type Handlers struct {
	Auth *auth.Handler
	User *user.Handler
	Post *post.Handler
}

func NewRouter(logger *slog.Logger, verifier middleware.TokenVerifier, h Handlers) *gin.Engine {
	// gin.SetMode(gin.ReleaseMode)
	r := gin.New()

	// 在 Use 之前註冊，所以 /health 不經過任何 middleware（gin 在註冊當下就把 handler chain 組好了）
	r.GET("/health", health)

	r.Use(
		middleware.RequestID(),
		middleware.AccessLog(logger),
		middleware.Recovery(logger),
	)
	// ---- 公開 ----
	h.Auth.RegisterRoutes(r)

	// ---- 有些公開，有些需登入 ----
	requireAuth := middleware.Auth(verifier)
	optionalAuth := middleware.OptionalAuth(verifier)

	h.Post.RegisterRoutes(r, requireAuth, optionalAuth)
	h.User.RegisterRoutes(r, requireAuth, optionalAuth)

	return r
}

func health(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"msg": "ok",
	})
}
