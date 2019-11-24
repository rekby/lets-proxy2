package profiler

import "net/http"

type secretHandler struct {
	argName string
	secret  string
	next    http.Handler
}

func (s secretHandler) ServeHTTP(resp http.ResponseWriter, req *http.Request) {
	querySecret := req.URL.Query().Get(s.argName)
	if s.secret == querySecret {
		s.next.ServeHTTP(resp, req)
		return
	}

	http.Error(resp, "Forbidden", http.StatusForbidden)
}
