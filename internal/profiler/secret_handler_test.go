package profiler

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gojuno/minimock/v3"
	"github.com/maxatome/go-testdeep"
)

//go:generate minimock -g -i net/http.Handler -o ./handler_mock_test.go

func TestSecretHandler_ServeHTTP(t *testing.T) {
	td := testdeep.NewT(t)
	mt := minimock.NewController(td)
	defer mt.Finish()

	nextHandler := NewHandlerMock(mt)
	nextCalled := false
	nextHandler.ServeHTTPMock.Set(func(resp http.ResponseWriter, req *http.Request) {
		resp.WriteHeader(http.StatusOK)
		_, _ = resp.Write([]byte("OK"))
		nextCalled = true
	})
	defer nextHandler.MinimockFinish()

	const argName = "pass"
	const secret = "asd"

	secretHandler := secretHandler{
		next:    nextHandler,
		argName: argName,
		secret:  secret,
	}

	nextCalled = false
	respWriter := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "http://test?pass=asd", nil)
	secretHandler.ServeHTTP(respWriter, req)
	respWriter.Flush()
	resp := respWriter.Result()
	td.Cmp(resp.StatusCode, http.StatusOK)
	td.True(nextCalled)

	nextCalled = false
	respWriter = httptest.NewRecorder()
	req = httptest.NewRequest(http.MethodGet, "http://test", nil)
	secretHandler.ServeHTTP(respWriter, req)
	respWriter.Flush()
	resp = respWriter.Result()
	td.Cmp(resp.StatusCode, http.StatusForbidden)
	td.False(nextCalled)
	nextCalled = false

	respWriter = httptest.NewRecorder()
	req = httptest.NewRequest(http.MethodGet, "http://test?pass=asdf", nil)
	secretHandler.ServeHTTP(respWriter, req)
	respWriter.Flush()
	resp = respWriter.Result()
	td.Cmp(resp.StatusCode, http.StatusForbidden)
	td.False(nextCalled)
}
