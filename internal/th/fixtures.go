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
	return fixenv.Cache(e, "",
		&fixenv.FixtureOptions{
			CleanupFunc: func() {
				server.Close()
			},
		}, func() (res string, err error) {
			server = httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
				u := *request.URL
				u.Scheme = ""
				u.Host = ""
				_, _ = writer.Write([]byte(u.String()))
			}))
			return server.URL, nil
		})
}

func TcpListener(e fixenv.Env) *net.TCPListener {
	var listener *net.TCPListener

	return fixenv.Cache(e, "",
		&fixenv.FixtureOptions{
			CleanupFunc: func() {
				_ = listener.Close()
			}},
		func() (res *net.TCPListener, err error) {
			listener, err = net.ListenTCP("tcp", &net.TCPAddr{IP: net.IPv4(127, 0, 0, 1)})
			return listener, err
		})
}

func TmpDir(e fixenv.Env) string {
	var dirPath string
	return fixenv.Cache(e, "", &fixenv.FixtureOptions{CleanupFunc: func() {
		_ = os.RemoveAll(dirPath)
	}}, func() (res string, err error) {
		dirPath, err = ioutil.TempDir("", "lets-proxy2-test-")
		return dirPath, err
	})
}
