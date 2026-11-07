package reqctx

import (
	"context"

	"github.com/google/uuid"
)

// 用未匯出的型別當 key，避免和其他套件撞 key
type requestIDKey struct{}
type userIDKey struct{}

func WithRequestID(ctx context.Context, id string) context.Context {
	return context.WithValue(ctx, requestIDKey{}, id)
}

func RequestID(ctx context.Context) string {
	id, _ := ctx.Value(requestIDKey{}).(string)
	return id
}

func WithUserID(ctx context.Context, id uuid.UUID) context.Context {
	return context.WithValue(ctx, userIDKey{}, id)
}

// 拿不到是正常情況（未登入 / 公開路由），所以回 bool，不回 error
func UserID(ctx context.Context) (uuid.UUID, bool) {
	id, ok := ctx.Value(userIDKey{}).(uuid.UUID)
	return id, ok && id != uuid.Nil
}
