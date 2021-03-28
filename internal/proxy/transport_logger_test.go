package proxy

import (
	"bytes"
	"context"
	"errors"
	"net/http"
	"net/url"
	"strings"
	"testing"
	"time"

	zc "github.com/rekby/zapcontext"

	"github.com/gojuno/minimock/v3"

	"github.com/maxatome/go-testdeep"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var _ http.RoundTripper = TransportLogger{}

//go:generate minimock -i net/http.RoundTripper -o http_round_tripper_mock_test.go -g
func TestTransportLogger(t *testing.T) {
	td := testdeep.NewT(t)

	buf := &bytes.Buffer{}

	encoder := zapcore.NewConsoleEncoder(zapcore.EncoderConfig{
		TimeKey:        "ts",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "caller",
		MessageKey:     "msg",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.LowercaseLevelEncoder,
		EncodeTime:     zapcore.ISO8601TimeEncoder,
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	})

	core := zapcore.NewCore(
		encoder,
		zap.CombineWriteSyncers(zapcore.AddSync(buf)),
		zapcore.DebugLevel,
	)
	logger := zap.New(core)

	ctx := zc.WithLogger(context.Background(), logger)

	mc := minimock.NewController(td)
	rtMock := NewRoundTripperMock(mc)

	resResp := &http.Response{StatusCode: 200}
	resErr := errors.New("resp-error")
	rtMock.RoundTripMock.Set(func(rp1 *http.Request) (rp2 *http.Response, err error) {
		time.Sleep(time.Millisecond)
		return resResp, resErr
	})

	tl := TransportLogger{rtMock}

	req := &http.Request{URL: &url.URL{}}
	req = req.WithContext(ctx)
	resp, err := tl.RoundTrip(req)
	td.Cmp(resp, resResp)
	td.Cmp(err, resErr)

	res := buf.String()
	if !strings.Contains(res, "status_code") {
		td.Error(res)
	}
}
