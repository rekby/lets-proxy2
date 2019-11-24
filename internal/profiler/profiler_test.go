package profiler

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/maxatome/go-testdeep"
	"go.uber.org/zap"
)

func TestNew(t *testing.T) {
	td := testdeep.NewT(t)
	profiler := New(zap.NewNop(), Config{AllowedNetworks: []string{"127.0.0.1/32", "::1/128"}})
	td.Len(profiler.secretHandler.AllowedNetworks, 2)
	td.NotNil(profiler.secretHandler.next)
	td.NotNil(profiler.secretHandler.logger)
}

func TestProfiler_ServeHTTP(t *testing.T) {
	td := testdeep.NewT(t)
	profiler := New(zap.NewNop(), Config{AllowedNetworks: []string{"127.0.0.1/32", "::1/128"}})

	respWriter := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "http://test/debug/pprof/", nil)
	req.RemoteAddr = "127.0.0.1:1234"
	profiler.ServeHTTP(respWriter, req)
	respWriter.Flush()
	resp := respWriter.Result()
	td.Cmp(resp.StatusCode, http.StatusOK)
	td.True(strings.Contains(resp.Header.Get("content-type"), "text/html"))
	_ = resp.Body.Close()
}
