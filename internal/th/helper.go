package th

import (
	"context"
	"io"

	"go.uber.org/zap/zaptest"

	zc "github.com/rekby/zapcontext"

	"go.uber.org/zap"
)

func TestContext(t zaptest.TestingT) (ctx context.Context, flush func()) {
	ctx, cancel := context.WithCancel(
		zc.WithLogger(context.Background(),
			Logger(t),
		),
	)
	flush = func() {
		cancel()
	}

	return ctx, flush
}

func NoLog(ctx context.Context) context.Context {
	return zc.WithLogger(ctx, zap.NewNop().WithOptions(zap.Development()))
}

func Logger(t zaptest.TestingT) *zap.Logger {
	return zaptest.NewLogger(t, zaptest.WrapOptions(zap.Development()))
}

func Close(closer io.Closer) {
	_ = closer.Close()
}
