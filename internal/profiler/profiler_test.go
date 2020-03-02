package profiler

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/rekby/lets-proxy2/internal/secrethandler"

	"github.com/maxatome/go-testdeep"
	"go.uber.org/zap"
)

func TestProfiler_ServeHTTP(t *testing.T) {
	td := testdeep.NewT(t)
	profiler := New(zap.NewNop(), Config{Config: secrethandler.Config{AllowEmptyPassword: true}})

	respWriter := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "http://test/debug/pprof/", nil)
	profiler.ServeHTTP(respWriter, req)
	respWriter.Flush()
	resp := respWriter.Result()
	td.Cmp(resp.StatusCode, http.StatusOK)
	td.True(strings.Contains(resp.Header.Get("content-type"), "text/html"))
	_ = resp.Body.Close()
}
