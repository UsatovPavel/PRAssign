package logging

import (
	"context"
	"log/slog"
)

type ctxKeyType struct{}

var ctxKey ctxKeyType

func WithLogger(ctx context.Context, l *slog.Logger) context.Context {
	return context.WithValue(ctx, ctxKey, l)
}

func FromContext(ctx context.Context) *slog.Logger {
	if v := ctx.Value(ctxKey); v != nil {
		if l, ok := v.(*slog.Logger); ok {
			return l
		}
	}
	return nil
}

func FromContextOrDefault(ctx context.Context, def *slog.Logger) *slog.Logger {
	l := FromContext(ctx)
	if l == nil {
		return def
	}
	return l
}
