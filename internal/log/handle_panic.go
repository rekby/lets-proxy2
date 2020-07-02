package log

import (
	"context"
	zc "github.com/rekby/zapcontext"
	"go.uber.org/zap"
)

// must called as defer HandlePanic(logger)
func HandlePanic(logger *zap.Logger) {
	err := recover()
	if err != nil {
		logger.DPanic("Panic handled", zap.Any("panic", err))
	}
}

// must called as defer HandlePanicCtx(ctx)
func HandlePanicCtx(ctx context.Context) {
	logger := zc.L(ctx)
	if logger == nil {
		logger, _ = zap.NewDevelopment()
		logger.Error("No context logger. Create tmp dev logger.")
	}
	HandlePanic(logger)
}
