package logger

import (
	"context"
	"log/slog"

	"github.com/ReilEgor/LinkShorteningService/internal/delivery/http/middleware"
)

type ContextHandler struct {
	slog.Handler
}

func (h *ContextHandler) Handle(ctx context.Context, r slog.Record) error {
	if reqID, ok := ctx.Value(middleware.RequestIDKey).(string); ok {
		r.AddAttrs(slog.String("request_id", reqID))
	}
	return h.Handler.Handle(ctx, r)
}
