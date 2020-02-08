package th

import (
	"context"

	zc "github.com/rekby/zapcontext"

	"go.uber.org/zap"
)

func TestContext() (ctx context.Context, flush func()) {
	ctx, cancel := context.WithCancel(zc.WithLogger(context.Background(), zap.NewNop().WithOptions(zap.Development())))
	flush = func() {
		cancel()
	}

	return ctx, flush
}

func NoLog(ctx context.Context) context.Context {
	return zc.WithLogger(ctx, zap.NewNop().WithOptions(zap.Development()))
}
