package th

import (
	"context"
	"crypto/tls"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"

	"github.com/rekby/fixenv"
	"go.uber.org/zap/zaptest"

	zc "github.com/rekby/zapcontext"

	"go.uber.org/zap"
)

func TestHttpServer(e fixenv.Env, handler http.Handler) *httptest.Server {
	server := httptest.NewServer(handler)
	e.T().Cleanup(server.Close)
	return server
}

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

func GetHttpClient() *http.Client {
	client := &http.Client{}
	client.Transport = &http.Transport{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: true,
		},
	}
	return client
}

func ErrorSubstringCmp(gotErr error, expect string) string {
	if expect == "" && gotErr == nil {
		return ""
	}

	var gotString string
	if gotErr != nil {
		gotString = gotErr.Error()
	}

	if expect == "" || !strings.Contains(gotString, expect) {
		return fmt.Sprintf("got string: '%v', expected: '%v'", gotString, expect)
	}
	return ""
}
