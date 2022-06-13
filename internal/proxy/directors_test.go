package proxy

import (
	"context"
	"net"
	"net/http"
	"testing"

	"github.com/rekby/lets-proxy2/internal/contextlabel"

	"github.com/gojuno/minimock/v3"

	"github.com/rekby/lets-proxy2/internal/th"

	"github.com/maxatome/go-testdeep"
)

func TestDirectorChain(t *testing.T) {
	mc := minimock.NewController(t)
	defer mc.Finish()

	var chain = NewDirectorChain()
	req := &http.Request{}
	chain.Director(req)

	d1 := NewDirectorMock(mc)
	d1.DirectorMock.Expect(req).Return(nil)
	d2 := NewDirectorMock(mc)
	d2.DirectorMock.Expect(req).Return(nil)
	chain = NewDirectorChain(d1, d2)
	chain.Director(req)
}

func TestDirectorDestMap(t *testing.T) {
	ctx, flush := th.TestContext(t)
	defer flush()

	td := testdeep.NewT(t)

	var req *http.Request

	m := map[string]string{
		(&net.TCPAddr{IP: net.ParseIP("1.2.3.1"), Port: 443}).String():        "1.1.1.1:80",
		(&net.TCPAddr{IP: net.ParseIP("1.2.3.2"), Port: 443}).String():        "2.2.2.2:80",
		(&net.TCPAddr{IP: net.ParseIP("::ffff:1.2.3.1"), Port: 443}).String(): "3.3.3.3:80",
	}

	d := NewDirectorDestMap(m)

	req = &http.Request{}
	req = req.WithContext(context.WithValue(
		ctx, http.LocalAddrContextKey, &net.TCPAddr{IP: net.ParseIP("aaa"), Port: 881}))
	d.Director(req)
	td.Nil(req.URL)

	req = &http.Request{}
	req = req.WithContext(context.WithValue(
		ctx, http.LocalAddrContextKey, &net.TCPAddr{IP: net.ParseIP("8.8.8.8"), Port: 443}))
	d.Director(req)
	td.Nil(req.URL)

	req = &http.Request{}
	req = req.WithContext(context.WithValue(
		ctx, http.LocalAddrContextKey, &net.TCPAddr{IP: net.ParseIP("1.2.3.1"), Port: 443}))
	d.Director(req)
	td.CmpDeeply(req.URL.Host, "3.3.3.3:80")

	req = &http.Request{}
	req = req.WithContext(context.WithValue(
		ctx, http.LocalAddrContextKey, &net.TCPAddr{IP: net.ParseIP("1.2.3.2"), Port: 443}))
	d.Director(req)
	td.CmpDeeply(req.URL.Host, "2.2.2.2:80")
}

func TestDirectorHost(t *testing.T) {
	td := testdeep.NewT(t)

	d := NewDirectorHost("haha:81")
	req := &http.Request{}
	d.Director(req)
	td.CmpDeeply(req.URL.Host, "haha:81")
}

func TestDirectorSameIP(t *testing.T) {
	ctx, flush := th.TestContext(t)
	defer flush()

	td := testdeep.NewT(t)

	d := NewDirectorSameIP(87)
	req := &http.Request{}
	req = req.WithContext(context.WithValue(
		ctx, http.LocalAddrContextKey, &net.TCPAddr{IP: net.ParseIP("1.2.3.4"), Port: 881}))
	d.Director(req)
	td.CmpDeeply(req.URL.Host, "1.2.3.4:87")
}

func TestDirectorSetHeaders(t *testing.T) {
	ctx, flush := th.TestContext(t)
	defer flush()

	td := testdeep.NewT(t)

	m := map[string]string{
		"TestConnectionID": "{{CONNECTION_ID}}",
		"TestIP":           "{{SOURCE_IP}}",
		"TestPort":         "{{SOURCE_PORT}}",
		"TestIPPort":       "{{SOURCE_IP}}:{{SOURCE_PORT}}",
		"TestVal":          "ddd",
		"TestProtocol":     "{{HTTP_PROTO}}",
	}

	d := NewDirectorSetHeaders(m)

	ctx = context.WithValue(ctx, contextlabel.ConnectionID, "123")

	req := &http.Request{RemoteAddr: "1.2.3.4:881"}
	req = req.WithContext(ctx)
	d.Director(req)
	td.CmpDeeply(req.Header.Get("TestConnectionID"), "123")
	td.CmpDeeply(req.Header.Get("TestIP"), "1.2.3.4")
	td.CmpDeeply(req.Header.Get("TestPort"), "881")
	td.CmpDeeply(req.Header.Get("TestIPPort"), "1.2.3.4:881")
	td.CmpDeeply(req.Header.Get("TestVal"), "ddd")
	td.CmpDeeply(req.Header.Get("TestProtocol"), "error protocol detection")

	req = &http.Request{RemoteAddr: "1.2.3.4:881"}
	ctx = context.WithValue(ctx, contextlabel.TLSConnection, true)
	req = req.WithContext(ctx)
	d.Director(req)
	td.CmpDeeply(req.Header.Get("TestProtocol"), "https")

	req = &http.Request{RemoteAddr: "1.2.3.4:881"}
	ctx = context.WithValue(ctx, contextlabel.TLSConnection, false)
	req = req.WithContext(ctx)
	d.Director(req)
	td.CmpDeeply(req.Header.Get("TestProtocol"), "http")
}
