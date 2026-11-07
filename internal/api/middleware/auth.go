package middleware

import (
	"errors"
	"net/http"
	"strings"

	"github.com/clary-work01/hot-post-api/internal/platform/httpx"
	"github.com/clary-work01/hot-post-api/internal/platform/reqctx"
	"github.com/clary-work01/hot-post-api/internal/platform/token"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type TokenVerifier interface {
	Verify(raw string) (uuid.UUID, error)
}

// Auth：一定要登入
func Auth(v TokenVerifier) gin.HandlerFunc {
	return authenticate(v, true)
}

// OptionalAuth：沒帶 token 當訪客；有帶就必須是有效的
func OptionalAuth(v TokenVerifier) gin.HandlerFunc {
	return authenticate(v, false)
}

func authenticate(v TokenVerifier, required bool) gin.HandlerFunc {
	return func(c *gin.Context) {
		header := c.GetHeader("Authorization")

		// 1. 完全沒帶：唯一一個兩種模式行為不同的地方
		if header == "" {
			if required {
				httpx.Error(c, http.StatusUnauthorized, httpx.ErrorBody{
					Code: "missing_token", Message: "未登入",
				})
				return
			}

			c.Next()
			return
		}

		// 2. 帶了但格式不對：client 有意要登入卻送錯，不能當訪客
		raw, ok := parseBearer(c.Request)
		if !ok {
			httpx.Error(c, http.StatusUnauthorized, httpx.ErrorBody{
				Code: "invalid_token", Message: "token無效",
			})
			return
		}

		// 3. 驗證
		userID, err := v.Verify(raw)
		if errors.Is(err, token.ErrExpiredToken) {
			httpx.Error(c, http.StatusUnauthorized, httpx.ErrorBody{
				Code: "token_expired", Message: "token過期",
			})
			return
		}
		if err != nil {
			httpx.Error(c, http.StatusUnauthorized, httpx.ErrorBody{
				Code: "invalid_token", Message: "token無效",
			})
			return
		}

		c.Request = c.Request.WithContext(reqctx.WithUserID(c.Request.Context(), userID))
		c.Next()
	}
}

func parseBearer(r *http.Request) (string, bool) {
	// Authorization header 的格式固定是：
	// Authorization: <scheme> <credentials>
	// Authorization: Bearer eyJhbGciOiJIUzI1NiIs...
	scheme, token, found := strings.Cut(r.Header.Get("Authorization"), " ")
	token = strings.TrimSpace(token)
	if !found || !strings.EqualFold(scheme, "Bearer") || token == "" {
		return "", false
	}
	return token, true
}
