package proxy

import (
	"net/http"
	"testing"

	"github.com/rekby/lets-proxy2/internal/th"

	"github.com/maxatome/go-testdeep"
)

func TestTransport_GetTransport(t *testing.T) {
	ctx, flush := th.TestContext(t)
	defer flush()

	td := testdeep.NewT(t)

	tr := Transport{}
	r, _ := http.NewRequest(http.MethodGet, "http://www.ru", nil)
	r = r.WithContext(ctx)
	httpTransport := tr.getTransport(r)
	td.True(httpTransport == defaultHTTPTransport) // equal pointers

	tr = Transport{IgnoreHTTPSCertificate: false}
	r, _ = http.NewRequest(http.MethodGet, "https://www.ru", nil)
	r = r.WithContext(ctx)
	httpTransport = tr.getTransport(r)
	td.True(httpTransport != defaultTransport()) // different pointers
	td.Cmp(httpTransport.TLSClientConfig.ServerName, "www.ru")
	td.Cmp(httpTransport.TLSClientConfig.InsecureSkipVerify, false)

	tr = Transport{IgnoreHTTPSCertificate: true}
	r, _ = http.NewRequest(http.MethodGet, "https://www.ru", nil)
	r = r.WithContext(ctx)
	httpTransport = tr.getTransport(r)
	td.True(httpTransport != defaultTransport()) // different pointers
	td.Cmp(httpTransport.TLSClientConfig.ServerName, "www.ru")
	td.Cmp(httpTransport.TLSClientConfig.InsecureSkipVerify, true)
}
