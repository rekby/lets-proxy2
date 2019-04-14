package th

import (
	"context"
	"fmt"
	"time"

	"github.com/rekby/zapcontext"

	"go.uber.org/zap"
)

func TestContext() (ctx context.Context, flush func()) {
	logger, err := zap.NewDevelopment()
	if err != nil {
		fmt.Print()
		panic(err)
	}

	ctx, cancel := context.WithCancel(zc.WithLogger(context.Background(), logger))
	flush = func() {
		cancel()
		time.Sleep(time.Millisecond)
		_ = logger.Sync()
	}
	return ctx, flush
}

func NoLog(ctx context.Context) context.Context {
	return zc.WithLogger(ctx, zap.NewNop().WithOptions(zap.Development()))
}
