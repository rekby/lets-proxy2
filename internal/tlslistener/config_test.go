package tlslistener

import (
	"net"
	"testing"

	"github.com/rekby/lets-proxy2/internal/th"

	"github.com/maxatome/go-testdeep"
)

func TestConfig_Apply(t *testing.T) {
	ctx, flush := th.TestContext(t)
	defer flush()

	td := testdeep.NewT(t)
	l := &ListenersHandler{}

	c := &Config{}
	err := c.Apply(ctx, l)
	td.CmpNoError(err)
	td.Empty(l.Listeners)
	td.Empty(l.ListenersForHandleTLS)

	c = &Config{
		TCPAddresses: []string{"asd"},
	}
	err = c.Apply(ctx, l)
	td.CmpError(err)
	td.Empty(l.Listeners)
	td.Empty(l.ListenersForHandleTLS)

	c = &Config{
		TLSAddresses: []string{"asd"},
	}
	err = c.Apply(ctx, l)
	td.CmpError(err)
	td.Empty(l.Listeners)
	td.Empty(l.ListenersForHandleTLS)

	const addr = "127.0.0.1"
	ports := getFreePorts(addr, 5)

	c = &Config{
		TCPAddresses: []string{addr + ":" + ports[0], addr + ":" + ports[1]},
		TLSAddresses: []string{addr + ":" + ports[2], addr + ":" + ports[3], addr + ":" + ports[4]},
	}
	err = c.Apply(ctx, l)

	defer func() {
		for _, listener := range l.ListenersForHandleTLS {
			_ = listener.Close()
		}
		for _, listener := range l.Listeners {
			_ = listener.Close()
		}
	}()

	td.CmpNoError(err)

	listenerAddresses := []string{l.Listeners[0].Addr().String(), l.Listeners[1].Addr().String()}
	td.CmpDeeply(listenerAddresses, []string{addr + ":" + ports[0], addr + ":" + ports[1]})

	tlsListenerAddresses := []string{l.ListenersForHandleTLS[0].Addr().String(), l.ListenersForHandleTLS[1].Addr().String(), l.ListenersForHandleTLS[2].Addr().String()}
	td.CmpDeeply(tlsListenerAddresses, []string{addr + ":" + ports[2], addr + ":" + ports[3], addr + ":" + ports[4]})
}

func getFreePorts(ip string, cnt int) []string {
	var res = make([]string, cnt)
	var listeners = make([]net.Listener, cnt)

	for i := 0; i < cnt; i++ {
		listener, err := net.Listen("tcp", ip+":0")
		if err != nil {
			panic(err)
		}
		_, port, err := net.SplitHostPort(listener.Addr().String())
		if err != nil {
			panic(err)
		}
		res[i] = port
		listeners[i] = listener
	}

	for _, l := range listeners {
		err := l.Close()
		if err != nil {
			panic(err)
		}
	}
	return res
}
