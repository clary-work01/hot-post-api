package logx

import (
	"context"
	"log/slog"
	"os"

	"github.com/clary-work01/hot-post-api/internal/platform/reqctx"
	"github.com/lmittmann/tint"
	"github.com/mattn/go-isatty"
)

// 自定義 contextHandler 實現 slog.Handler 接口
type contextHandler struct {
	slog.Handler
}

// 呼叫 XxxContext 的方法 , 自動幫每條 log 加上 request_id
func (h *contextHandler) Handle(ctx context.Context, r slog.Record) error {
	if rid := reqctx.RequestID(ctx); rid != "" {
		r.AddAttrs(slog.String("request_id", rid))
	}

	if uid, ok := reqctx.UserID(ctx); ok {
		r.AddAttrs(slog.String("user_id", uid.String()))
	}
	return h.Handler.Handle(ctx, r)
}

func (h *contextHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return &contextHandler{h.Handler.WithAttrs(attrs)}
}

func (h *contextHandler) WithGroup(name string) slog.Handler {
	return &contextHandler{h.Handler.WithGroup(name)}
}

func New(isProd bool) *slog.Logger {
	var h slog.Handler

	if isProd {
		h = slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo})
	} else {
		// h = slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug})
		h = tint.NewTextHandler(os.Stdout, &tint.Options{
			Level:      slog.LevelDebug,
			TimeFormat: "15:04:05.000",
			// 輸出不是 terminal（例如 > app.log 或 | grep）時關掉顏色，
			// 否則檔案裡會塞滿 \x1b[32m 這種 ANSI 控制碼
			NoColor: !isatty.IsTerminal(os.Stdout.Fd()),
		})
	}

	return slog.New(&contextHandler{h})
}
