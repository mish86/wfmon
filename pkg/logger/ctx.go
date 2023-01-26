package log

import (
	"context"

	"go.uber.org/zap"
)

type ctxKey struct{}

// Returns logger from context if exists, otherwise global.
func CtxL(ctx context.Context) *zap.Logger {
	if l, ok := ctx.Value(ctxKey{}).(*zap.Logger); ok {
		return l
	}
	return zap.L()
}

// Returns sugared logger from context if exists, otherwise global.
func CtxS(ctx context.Context) *zap.SugaredLogger {
	if l, ok := ctx.Value(ctxKey{}).(*zap.Logger); ok {
		return l.Sugar()
	}
	return zap.S()
}

// Adds logger to context.
func CtxWithLog(ctx context.Context, l *zap.Logger) context.Context {
	return context.WithValue(ctx, ctxKey{}, l)
}
