package secrethandler

import (
	"net"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/rekby/lets-proxy2/internal/th"

	"go.uber.org/zap/zaptest"

	"github.com/gojuno/minimock/v3"
	"github.com/maxatome/go-testdeep"
	"go.uber.org/zap"
)

//go:generate minimock -g -i net/http.Handler -o ./handler_mock_test.go

func TestNew(t *testing.T) {
	td := testdeep.NewT(t)
	mux := http.ServeMux{}
	h := New(th.Logger(td), Config{AllowedNetworks: []string{"127.0.0.1/32", "::1/128"}}, &mux)
	td.Len(h.allowedNetworks, 2)
	td.NotNil(h.next)
	td.NotNil(h.logger)
	td.Cmp(h.next, &mux)
}

func TestSecretHandler_TestSources(t *testing.T) {
	td := testdeep.NewT(t)
	mt := minimock.NewController(td)
	defer mt.Finish()

	logger := zaptest.NewLogger(td, zaptest.WrapOptions(zap.Development()))

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

	// Test source
	secretHandler := SecretHandler{
		next:               nextHandler,
		allowedNetworks:    localNetworks,
		logger:             logger,
		allowEmptyPassword: true,
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
	_ = resp.Body.Close()

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
	_ = resp.Body.Close()

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
	_ = resp.Body.Close()
}

func TestPassword(t *testing.T) {
	td := testdeep.NewT(t)
	mt := minimock.NewController(td)
	defer mt.Finish()

	logger := zaptest.NewLogger(td, zaptest.WrapOptions(zap.Development()))

	nextHandler := NewHandlerMock(mt)
	nextCalled := false
	nextHandler.ServeHTTPMock.Set(func(resp http.ResponseWriter, req *http.Request) {
		resp.WriteHeader(http.StatusOK)
		_, _ = resp.Write([]byte("OK"))
		nextCalled = true
	})
	defer nextHandler.MinimockFinish()

	const password = "123"
	secretHandler := SecretHandler{
		next:               nextHandler,
		logger:             logger,
		allowEmptyPassword: false,
		password:           password,
	}

	nextCalled = false
	respWriter := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "http://test?password="+password, nil)
	secretHandler.ServeHTTP(respWriter, req)
	respWriter.Flush()
	resp := respWriter.Result()
	td.Cmp(resp.StatusCode, http.StatusOK)
	td.True(nextCalled)
	_ = resp.Body.Close()

	nextCalled = false
	respWriter = httptest.NewRecorder()
	req = httptest.NewRequest(http.MethodGet, "http://test", nil)
	secretHandler.ServeHTTP(respWriter, req)
	respWriter.Flush()
	resp = respWriter.Result()
	td.Cmp(resp.StatusCode, http.StatusForbidden)
	td.False(nextCalled)
	nextCalled = false
	_ = resp.Body.Close()
}
