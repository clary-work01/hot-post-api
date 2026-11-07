package middleware

import (
	"errors"
	"log/slog"
	"net"
	"net/http"
	"os"
	"runtime/debug"
	"strings"

	"github.com/gin-gonic/gin"
)

func Recovery(logger *slog.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			rec := recover()
			if rec == nil { // 沒有 panic , 正常結束
				return
			}

			// 1. http.ErrAbortHandler 是刻意中止請求，照慣例重新 panic 交給 net/http 處理
			if err, ok := rec.(error); ok && errors.Is(err, http.ErrAbortHandler) {
				panic(err)
			}

			// 2. client 斷線 (broken pipe / connection reset)：寫不回去了，記個 warn 就好
			if isBrokenPipe(rec) {
				logger.LogAttrs(c.Request.Context(), slog.LevelWarn, "client disconnected",
					slog.Any("error", rec),
					slog.String("path", c.Request.URL.Path),
				)
				c.Error(rec.(error))
				c.Abort()
				return
			}

			// 3. 真正的 panic：記完整 stack
			stack := debug.Stack()
			logger.ErrorContext(c.Request.Context(), "panic recovered",
				slog.Any("panic", rec),
				slog.String("method", c.Request.Method),
				slog.String("path", c.Request.URL.Path),
				slog.String("stack", string(stack)),
			)

			// 只在 dev：多印一份有換行的原始 stack 給人看
			if gin.Mode() == gin.DebugMode {
				os.Stderr.Write(stack)
			}

			// 4. 已經寫過 response header 就不能再寫，只能 abort
			if c.Writer.Written() {
				c.Abort()
				return
			}

			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"code":    "INTERNAL_ERROR",
				"message": "internal server error",
			})
		}()

		c.Next()
	}
}

func isBrokenPipe(rec any) bool {
	err, ok := rec.(error)
	if !ok {
		return false
	}

	var opErr *net.OpError
	if errors.As(err, &opErr) {
		var sysErr *os.SyscallError
		if errors.As(opErr.Err, &sysErr) {
			msg := strings.ToLower(sysErr.Error())
			return strings.Contains(msg, "broken pipe") ||
				strings.Contains(msg, "connection reset by peer")
		}
	}
	return false
}
