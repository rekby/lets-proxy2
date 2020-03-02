package dns

import (
	"net"
	"testing"

	"github.com/pkg/errors"

	"github.com/gojuno/minimock/v3"

	"github.com/maxatome/go-testdeep"
	"github.com/rekby/lets-proxy2/internal/th"
)

var (
	_ ResolverInterface = Parallel{}
)

func TestParallel(t *testing.T) {
	ctx, cancel := th.TestContext(t)
	defer cancel()

	td := testdeep.NewT(t)
	mc := minimock.NewController(td)

	var ips []net.IPAddr
	var err error

	p := NewParallel()
	ips, err = p.LookupIPAddr(ctx, "123")
	td.CmpNoError(err)
	td.Nil(ips)

	r1 := NewResolverInterfaceMock(mc)
	r2 := NewResolverInterfaceMock(mc)

	p = NewParallel(r1)

	r1.LookupIPAddrMock.When(ctx, "1").Then([]net.IPAddr{{IP: net.ParseIP("1.2.3.4")}}, nil)
	ips, err = p.LookupIPAddr(ctx, "1")
	td.CmpNoError(err)
	td.CmpDeeply(ips, []net.IPAddr{{IP: net.ParseIP("1.2.3.4")}})

	testErr := errors.New("test2")
	r1.LookupIPAddrMock.When(ctx, "2").Then(nil, testErr)
	ips, err = p.LookupIPAddr(ctx, "2")
	td.CmpDeeply(err, testErr)
	td.Nil(ips)

	p = NewParallel(r1, r2)
	r1.LookupIPAddrMock.When(ctx, "3").Then([]net.IPAddr{{IP: net.ParseIP("1.2.3.4")}}, nil)
	r2.LookupIPAddrMock.When(ctx, "3").Then([]net.IPAddr{{IP: net.ParseIP("4.5.6.7")}}, nil)
	ips, err = p.LookupIPAddr(ctx, "3")
	td.CmpNoError(err)
	td.CmpDeeply(ips, []net.IPAddr{{IP: net.ParseIP("1.2.3.4")}, {IP: net.ParseIP("4.5.6.7")}})

	r1.LookupIPAddrMock.When(ctx, "4").Then([]net.IPAddr{{IP: net.ParseIP("1.2.3.4")}}, nil)
	r2.LookupIPAddrMock.When(ctx, "4").Then(nil, errors.New("test4"))
	ips, err = p.LookupIPAddr(ctx, "4")
	td.CmpNoError(err)
	td.CmpDeeply(ips, []net.IPAddr{{IP: net.ParseIP("1.2.3.4")}})

	r1.LookupIPAddrMock.When(ctx, "5").Then(nil, errors.New("test5"))
	r2.LookupIPAddrMock.When(ctx, "5").Then([]net.IPAddr{{IP: net.ParseIP("4.5.6.7")}}, nil)
	ips, err = p.LookupIPAddr(ctx, "5")
	td.CmpNoError(err)
	td.CmpDeeply(ips, []net.IPAddr{{IP: net.ParseIP("4.5.6.7")}})

	error61 := errors.New("test6-1")
	error62 := errors.New("test6-2")
	r1.LookupIPAddrMock.When(ctx, "6").Then(nil, error61)
	r2.LookupIPAddrMock.When(ctx, "6").Then(nil, error62)
	ips, err = p.LookupIPAddr(ctx, "6")
	td.Any(err, []interface{}{error61, error62})
	td.Nil(ips)
}

func TestParallelReadl(t *testing.T) {
	ctx, cancel := th.TestContext(t)
	defer cancel()

	td := testdeep.NewT(t)

	r := NewParallel(NewResolver("8.8.8.8:53"), NewResolver("4.4.4.4:53"))
	ips, err := r.LookupIPAddr(ctx, "one.one.one.one")
	td.CmpNoError(err)
	td.Contains(ips,
		testdeep.Any(
			net.IPAddr{IP: net.IPv4(1, 1, 1, 1).To4()},
			net.IPAddr{IP: net.IPv4(1, 1, 1, 1).To16()},
		),
	)
}
