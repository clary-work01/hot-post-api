package api

import (
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"

	"github.com/clary-work01/hot-post-api/internal/domain"
	"github.com/clary-work01/hot-post-api/internal/domain/user"
)

func decodeJSON(w http.ResponseWriter, r *http.Request, data any) error {
	maxBytes := 1_048_576
	r.Body = http.MaxBytesReader(w, r.Body, int64(maxBytes))

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	return decoder.Decode(data)
}

type jsonResp struct {
	Message string `json:"message"`
	Data    any    `json:"data,omitempty"`
}

func JSON(w http.ResponseWriter, status int, msg string, data any) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	return json.NewEncoder(w).Encode(jsonResp{
		Message: msg,
		Data:    data,
	})
}

func writeError(w http.ResponseWriter, r *http.Request, err error) {
	switch {
	case errors.Is(err, user.ErrCannotFollowSelf):
		JSON(w, http.StatusBadRequest, "不能關注自己", nil) // 400

	case errors.Is(err, user.ErrUsernameTaken):
		JSON(w, http.StatusBadRequest, "用戶名已存在", nil) // 400

	case errors.Is(err, domain.ErrInvalidInput):
		JSON(w, http.StatusBadRequest, "參數錯誤", nil) // 400

	case errors.Is(err, domain.ErrNotFound):
		JSON(w, http.StatusNotFound, "資源不存在", nil) // 404

	case errors.Is(err, domain.ErrConflict):
		JSON(w, http.StatusConflict, "已存在", nil) // 409

	default:
		JSON(w, http.StatusInternalServerError, "系統錯誤", nil)

		// 只有非預期錯誤才 log 完整 error chain, 4xx 不用 log 成 error
		slog.Error("internal server error", "method", r.Method, "path", r.URL.Path, "error", err)
	}
}
