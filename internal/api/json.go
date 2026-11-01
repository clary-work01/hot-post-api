package api

import (
	"encoding/json"
	"errors"

	"log/slog"
	"net/http"

	"github.com/clary-work01/hot-post-api/internal/domain"
	"github.com/clary-work01/hot-post-api/internal/domain/auth"
	"github.com/clary-work01/hot-post-api/internal/domain/user"
	"github.com/clary-work01/hot-post-api/internal/token"
)

// ---- request ----

func decodeJSON(w http.ResponseWriter, r *http.Request, data any) error {
	maxBytes := 1_048_576
	r.Body = http.MaxBytesReader(w, r.Body, int64(maxBytes))

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	return decoder.Decode(data)
}

// ---- response ----

type errorCode string

const (
	// auth
	codeInvalidCredentials  errorCode = "invalid_credentials"
	codeInvalidRefreshToken errorCode = "invalid_refresh_token"

	// jwt
	codeInvalidToken errorCode = "invalid_token"
	codeTokenExpired errorCode = "token_expired"

	// user
	codeUsernameTaken    errorCode = "username_taken"
	codeCannotFollowSelf errorCode = "cannot_follow_self"

	// 共用
	codeInvalidInput errorCode = "invalid_input"
	codeUnauthorized errorCode = "unauthorized"
	codeNotFound     errorCode = "not_found"
	codeConflict     errorCode = "conflict"
	codeInternal     errorCode = "internal_error"
)

type jsonResp struct {
	Code    errorCode `json:"code,omitempty"` // 只有錯誤回應會有
	Message string    `json:"message"`
	Data    any       `json:"data,omitempty"`
}

func encodeJSON(w http.ResponseWriter, status int, resp jsonResp) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	if err := json.NewEncoder(w).Encode(resp); err != nil {
		slog.Error("encode json response", "error", err)
	}
}

// 成功回應
func JSON(w http.ResponseWriter, status int, msg string, data any) {
	encodeJSON(w, status, jsonResp{
		Message: msg,
		Data:    data,
	})
}

// 錯誤回應
func errorJSON(w http.ResponseWriter, status int, code errorCode, msg string) {
	encodeJSON(w, status, jsonResp{
		Code:    code,
		Message: msg,
	})
}

func writeError(w http.ResponseWriter, r *http.Request, err error) {
	switch {
	// ---- auth ----

	case errors.Is(err, auth.ErrInvalidCredentials):
		errorJSON(w, http.StatusUnauthorized, codeInvalidCredentials, "帳密錯誤") // 401

	case errors.Is(err, auth.ErrInvalidRefreshToken):
		errorJSON(w, http.StatusUnauthorized, codeInvalidRefreshToken, "請重新登入") // 401

	// ---- jwt ----
	case errors.Is(err, token.ErrInvalidToken):
		errorJSON(w, http.StatusUnauthorized, codeInvalidToken, "token無效") // 401

	case errors.Is(err, token.ErrExpiredToken):
		errorJSON(w, http.StatusUnauthorized, codeTokenExpired, "token過期") // 401

	// ---- user ----
	case errors.Is(err, user.ErrCannotFollowSelf):
		errorJSON(w, http.StatusBadRequest, codeCannotFollowSelf, "不能關注自己") // 400

	case errors.Is(err, user.ErrUsernameTaken):
		errorJSON(w, http.StatusConflict, codeUsernameTaken, "用戶名已存在") // 409

	// ---- 共用 ----
	case errors.Is(err, domain.ErrInvalidInput):
		errorJSON(w, http.StatusBadRequest, codeInvalidInput, "參數錯誤") // 400

	case errors.Is(err, domain.ErrUnauthorized):
		errorJSON(w, http.StatusUnauthorized, codeUnauthorized, "未授權") // 401

	case errors.Is(err, domain.ErrNotFound):
		errorJSON(w, http.StatusNotFound, codeNotFound, "資源不存在") // 404

	case errors.Is(err, domain.ErrConflict):
		errorJSON(w, http.StatusConflict, codeConflict, "已存在") // 409

	default:
		errorJSON(w, http.StatusInternalServerError, codeInternal, "系統錯誤")

		// 只有非預期錯誤才 log 完整 error chain, 4xx 不用 log 成 error
		slog.Error("internal server error", "method", r.Method, "path", r.URL.Path, "error", err)
	}
}
