package zc

import (
	"context"
	"go.uber.org/zap"
)

type loggerContextMark struct{}

// Return copy of ctx with assiciated logger
func WithLogger(ctx context.Context, l *zap.Logger) context.Context {
	return context.WithValue(ctx, loggerContextMark{}, l)
}

func WithSugarLogger(ctx context.Context, l *zap.SugaredLogger) context.Context {
	return context.WithValue(ctx, loggerContextMark{}, l)
}

// Return associated logger. Nop if none.
func L(ctx context.Context) *zap.Logger {
	val := ctx.Value(loggerContextMark{})
	if logger, ok := val.(*zap.Logger); ok {
		return logger
	}
	if logger, ok := val.(*zap.SugaredLogger); ok {
		return logger.Desugar()
	}
	return nil
}

func LNop(ctx context.Context) (res *zap.Logger) {
	res = L(ctx)
	if res == nil {
		res = zap.NewNop()
	}
	return res
}

func S(ctx context.Context) *zap.SugaredLogger {
	res := ctx.Value(loggerContextMark{})
	if logger, ok := res.(*zap.SugaredLogger); ok {
		return logger
	}
	if logger, ok := res.(*zap.Logger); ok {
		return logger.Sugar()
	}
	return nil
}

func SNop(ctx context.Context) (res *zap.SugaredLogger) {
	res = S(ctx)
	if res == nil {
		res = zap.NewNop().Sugar()
	}
	return res
}
