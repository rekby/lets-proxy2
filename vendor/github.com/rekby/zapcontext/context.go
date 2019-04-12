package zc

import (
	"context"
	"go.uber.org/zap"
)

type loggerContextMark struct{}

var (
	defaultLogger        *zap.Logger
	defaultSugaredLogger *zap.SugaredLogger
	nopLogger            = zap.NewNop()
	nopSugar             = zap.NewNop().Sugar()
)

// SetDefaultLogger set logger, which will return with L and S function if no logger in context.
// It is NOT thread safe and no syncronization with L and S
func SetDefaultLogger(logger *zap.Logger) {
	defaultLogger = logger
	if logger == nil {
		defaultSugaredLogger = nil
	} else {
		defaultSugaredLogger = logger.Sugar()
	}
}

// Return copy of ctx with associated logger
func WithLogger(ctx context.Context, l *zap.Logger) context.Context {
	return context.WithValue(ctx, loggerContextMark{}, l)
}

// Return copy of ctx with associated logger
func WithSugarLogger(ctx context.Context, l *zap.SugaredLogger) context.Context {
	return context.WithValue(ctx, loggerContextMark{}, l)
}

// Return associated logger. Default logger if no ctx logger yet. (default of default is nil)
func L(ctx context.Context) *zap.Logger {
	if ctx == nil {
		return defaultLogger
	}

	val := ctx.Value(loggerContextMark{})
	if logger, ok := val.(*zap.Logger); ok {
		return logger
	}
	if logger, ok := val.(*zap.SugaredLogger); ok {
		return logger.Desugar()
	}
	return defaultLogger
}

// Return associated logger. Nop if no ctx logger yet.
func LNop(ctx context.Context) (res *zap.Logger) {
	if ctx == nil {
		return nopLogger
	}

	val := ctx.Value(loggerContextMark{})
	if logger, ok := val.(*zap.Logger); ok {
		return logger
	}
	if logger, ok := val.(*zap.SugaredLogger); ok {
		return logger.Desugar()
	}
	return nopLogger
}

// Return associated logger. Default logger if no ctx logger yet. (default of default is nil)
func S(ctx context.Context) *zap.SugaredLogger {
	if ctx == nil {
		return defaultSugaredLogger
	}

	res := ctx.Value(loggerContextMark{})
	if logger, ok := res.(*zap.SugaredLogger); ok {
		return logger
	}
	if logger, ok := res.(*zap.Logger); ok {
		return logger.Sugar()
	}
	return defaultSugaredLogger
}

// Return associated logger. Nop if no ctx logger yet.
func SNop(ctx context.Context) *zap.SugaredLogger {
	if ctx == nil {
		return nopSugar
	}

	res := ctx.Value(loggerContextMark{})
	if logger, ok := res.(*zap.SugaredLogger); ok {
		return logger
	}
	if logger, ok := res.(*zap.Logger); ok {
		return logger.Sugar()
	}
	return nopSugar
}
