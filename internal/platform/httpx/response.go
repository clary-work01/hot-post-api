package httpx

import (
	"context"
	"errors"
	"log/slog"

	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

const CtxKeyErrCode = "err_code"

type ErrorBody struct {
	Code    string            `json:"code"`
	Message string            `json:"message"`
	Details map[string]string `json:"details,omitempty"`
}

// 成功回應
func JSON(c *gin.Context, status int, data any) {
	c.JSON(status, gin.H{"data": data})
}

// 成功但沒有內容（PUT / DELETE 等）
func NoContent(c *gin.Context) {
	c.Status(http.StatusNoContent)
}

// 所有錯誤回應都走這裡，確保格式一致
func Error(c *gin.Context, status int, body ErrorBody) {
	c.Set(CtxKeyErrCode, body.Code)
	c.Abort()
	c.JSON(status, gin.H{"error": body})
}

func InvalidInputError(c *gin.Context, err error) {
	var mbe *http.MaxBytesError
	if errors.As(err, &mbe) {
		Error(c, http.StatusRequestEntityTooLarge, ErrorBody{Code: "body_too_large", Message: "請求內容過大"})
		return
	}

	body := ErrorBody{Code: "invalid_input", Message: "參數錯誤"}
	var ve validator.ValidationErrors
	if errors.As(err, &ve) {
		body.Details = make(map[string]string, len(ve))
		for _, fe := range ve {
			body.Details[fe.Field()] = fe.Tag()
		}
	}
	Error(c, http.StatusBadRequest, body)
}

// 非預期錯誤：這是邊界，唯一 log 一次完整的 error chain
func ServerError(c *gin.Context, err error) {
	ctx := c.Request.Context()

	// client 自己斷線，不算我們的錯，不記 error
	if errors.Is(err, context.Canceled) && ctx.Err() != nil {
		c.Set(CtxKeyErrCode, "client_closed")
		c.AbortWithStatus(499)
		return
	}

	slog.ErrorContext(ctx, "internal server error",
		"method", c.Request.Method,
		"path", c.FullPath(),
		"error", err,
	)

	Error(c, http.StatusInternalServerError, ErrorBody{Code: "internal_error", Message: "系統錯誤"})
}
