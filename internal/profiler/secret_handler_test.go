package profiler

import (
	"net"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gojuno/minimock/v3"
	"github.com/maxatome/go-testdeep"
	"go.uber.org/zap"
)

//go:generate minimock -g -i net/http.Handler -o ./handler_mock_test.go

func TestSecretHandler_ServeHTTP(t *testing.T) {
	td := testdeep.NewT(t)
	mt := minimock.NewController(td)
	defer mt.Finish()

	_, n, _ := net.ParseCIDR("1.2.3.4/32")
	var localNetworks = []net.IPNet{*n}

	nextHandler := NewHandlerMock(mt)
	nextCalled := false
	nextHandler.ServeHTTPMock.Set(func(resp http.ResponseWriter, req *http.Request) {
		resp.WriteHeader(http.StatusOK)
		_, _ = resp.Write([]byte("OK"))
		nextCalled = true
	})
	defer nextHandler.MinimockFinish()

	secretHandler := secretHandler{
		next:            nextHandler,
		AllowedNetworks: localNetworks,
		logger:          zap.NewNop(),
	}

	nextCalled = false
	respWriter := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "http://test", nil)
	req.RemoteAddr = "1.2.3.4:1234"
	secretHandler.ServeHTTP(respWriter, req)
	respWriter.Flush()
	resp := respWriter.Result()
	td.Cmp(resp.StatusCode, http.StatusOK)
	td.True(nextCalled)

	nextCalled = false
	respWriter = httptest.NewRecorder()
	req = httptest.NewRequest(http.MethodGet, "http://test", nil)
	req.RemoteAddr = "2.3.4.5:1234"
	secretHandler.ServeHTTP(respWriter, req)
	respWriter.Flush()
	resp = respWriter.Result()
	td.Cmp(resp.StatusCode, http.StatusForbidden)
	td.False(nextCalled)
	nextCalled = false

	nextCalled = false
	respWriter = httptest.NewRecorder()
	req = httptest.NewRequest(http.MethodGet, "http://test", nil)
	req.RemoteAddr = "asdf"
	secretHandler.ServeHTTP(respWriter, req)
	respWriter.Flush()
	resp = respWriter.Result()
	td.Cmp(resp.StatusCode, http.StatusInternalServerError)
	td.False(nextCalled)
	nextCalled = false
}
