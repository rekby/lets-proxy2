package th

import (
	"io/ioutil"
	"net"
	"net/http"
	"net/http/httptest"
	"os"

	"github.com/rekby/fixenv"
)

func HttpQuery(e fixenv.Env) string {
	var server *httptest.Server
	return e.Cache(nil,
		&fixenv.FixtureOptions{
			CleanupFunc: func() {
				server.Close()
			},
		}, func() (res interface{}, err error) {
			server = httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
				u := *request.URL
				u.Scheme = ""
				u.Host = ""
				_, _ = writer.Write([]byte(u.String()))
			}))
			return server.URL, nil
		}).(string)
}

func TcpListener(e fixenv.Env) *net.TCPListener {
	var listener *net.TCPListener

	return e.Cache(nil,
		&fixenv.FixtureOptions{
			CleanupFunc: func() {
				_ = listener.Close()
			}},
		func() (res interface{}, err error) {
			listener, err = net.ListenTCP("tcp", &net.TCPAddr{IP: net.IPv4(127, 0, 0, 1)})
			return listener, err
		}).(*net.TCPListener)
}

func TmpDir(e fixenv.Env) string {
	var dirPath string
	return e.Cache(nil, &fixenv.FixtureOptions{CleanupFunc: func() {
		_ = os.RemoveAll(dirPath)
	}}, func() (res interface{}, err error) {
		dirPath, err = ioutil.TempDir("", "lets-proxy2-test-")
		return dirPath, err
	}).(string)
}
