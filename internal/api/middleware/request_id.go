package middleware

import (
	"regexp"

	"github.com/clary-work01/hot-post-api/internal/platform/reqctx"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// 外部傳進來的值不能完全信任：限制字元與長度，防止 log injection
var validRID = regexp.MustCompile(`^[a-zA-Z0-9\-_]{1,64}$`)

func RequestID() gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.GetHeader("X-Request-ID")

		if !validRID.MatchString(id) {
			id = uuid.NewString()
		}

		c.Request = c.Request.WithContext(reqctx.WithRequestID(c.Request.Context(), id))
		c.Header("X-Request-ID", id) // 回傳給前端，使用者回報問題時可以附上

		c.Next()
	}
}
