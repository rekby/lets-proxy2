package domain_checker

import (
	"net"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/maxatome/go-testdeep"
	"github.com/rekby/lets-proxy2/internal/th"
)

func TestAwsPublicIPs(t *testing.T) {
	ctx, cancel := th.TestContext(t)
	defer cancel()
	td := testdeep.NewT(t)

	getMetadata := func(p string) (string, error) {
		switch p {
		case "public-ipv4":
			return "1.2.3.4", nil
		case "network/interfaces/macs/":
			return "d0:0d:aa:1b:c4:ba", nil
		case "network/interfaces/macs/d0:0d:aa:1b:c4:ba/ipv6s":
			return "2a02::3", nil
		default:
			panic(p)
		}
	}
	getIPsFunc := awsPublicIPs(getMetadata)
	ips, err := getIPsFunc(ctx)
	td.CmpNoError(err)
	td.CmpDeeply(ips, []net.IP{net.ParseIP("1.2.3.4"), net.ParseIP("2a02::3")})
}

func TestGetIpByExternalRequest(t *testing.T) {
	ctx, cancel := th.TestContext(t)
	defer cancel()

	td := testdeep.NewT(t)

	mux := http.ServeMux{}
	mux.HandleFunc("/", func(writer http.ResponseWriter, request *http.Request) {
		_, _ = writer.Write([]byte("1.2.3.4"))
	})
	server := httptest.NewServer(&mux)
	defer server.Close()

	ips, err := getIPByExternalRequest(ctx, server.URL)
	td.CmpNoError(err)
	if len(ips) == 1 {
		// ipv4 only host or ipv6 only host
		td.Cmp(ips[0], net.ParseIP("1.2.3.4"))
	} else {
		// test server answer same ip twice
		td.Cmp(ips[0], net.ParseIP("1.2.3.4"))
		td.Cmp(ips[1], net.ParseIP("1.2.3.4"))
	}
}

func TestCreateGetSelfPublicBinded(t *testing.T) {
	ctx, cancel := th.TestContext(t)
	defer cancel()

	td := testdeep.NewT(t)

	var binded InterfacesAddrFunc = func() (addrs []net.Addr, e error) {
		return []net.Addr{
			&net.IPNet{IP: net.ParseIP("1.2.3.4"), Mask: net.CIDRMask(32, 32)},
			&net.IPNet{IP: net.ParseIP("127.0.0.1"), Mask: net.CIDRMask(32, 32)},
		}, nil
	}

	f := SelfBindedPublicIPs(binded)
	ips, err := f(ctx)
	td.CmpDeeply(len(ips), 1)
	td.True(ips[0].Equal(net.ParseIP("1.2.3.4")))
	td.CmpNoError(err)
}
