package middleware

import (
	"log/slog"
	"time"

	"github.com/clary-work01/hot-post-api/internal/platform/httpx"
	"github.com/gin-gonic/gin"
)

// access log：只負責紀錄，固定 INFO
func AccessLog(logger *slog.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.Request.URL.Path == "/health" { // 健康檢查會被頻繁呼叫，跳過
			c.Next()
			return
		}

		start := time.Now()
		c.Next()

		logger.LogAttrs(c.Request.Context(), slog.LevelInfo, "http request",
			slog.String("method", c.Request.Method),
			slog.String("route", c.FullPath()),
			slog.String("path", c.Request.URL.Path),
			slog.Int("status", c.Writer.Status()),
			slog.Int64("latency_ms", time.Since(start).Milliseconds()),
			slog.String("client_ip", c.ClientIP()),
			slog.String("remote_ip", c.RemoteIP()),
			slog.String("err_code", c.GetString(httpx.CtxKeyErrCode)), // 沒錯時為空字串
		)
	}
}
