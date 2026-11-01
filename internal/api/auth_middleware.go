package api

import (
	"context"
	"net/http"
	"strings"

	"github.com/clary-work01/hot-post-api/internal/domain"
	"github.com/google/uuid"
)

type ctxKey struct{} // 自訂未匯出型別，避免和別的 package 的 context key 撞名

type TokenVerifier interface {
	Verify(raw string) (uuid.UUID, error)
}

func RequireAuth(v TokenVerifier) Middleware {
	return func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			raw, ok := bearerToken(r)
			if !ok {
				writeError(w, r, domain.ErrUnauthorized)
				return
			}

			userID, err := v.Verify(raw)
			if err != nil {
				writeError(w, r, err)
				return
			}

			// 做一個「多帶了 userID」的新 context
			ctx := context.WithValue(r.Context(), ctxKey{}, userID)
			// 換上新 context，交給下一個 handler
			next(w, r.WithContext(ctx))
		}
	}
}

// 只負責「拿」，拿不到是正常情況
func userIDFromCtx(ctx context.Context) (uuid.UUID, bool) {
	id, ok := ctx.Value(ctxKey{}).(uuid.UUID)
	return id, ok
}

// 從 context 拿登入者 id
// 一定要有，沒有就回 401
func currentUser(w http.ResponseWriter, r *http.Request) (uuid.UUID, bool) {
	id, ok := userIDFromCtx(r.Context())
	if !ok {
		writeError(w, r, domain.ErrUnauthorized)
	}
	return id, ok
}

func bearerToken(r *http.Request) (string, bool) {
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
