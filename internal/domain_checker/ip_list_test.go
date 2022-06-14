//nolint:golint
package domain_checker

import (
	"context"
	"errors"
	"net"
	"testing"
	"time"

	zc "github.com/rekby/zapcontext"
	"go.uber.org/zap"

	"github.com/rekby/lets-proxy2/internal/th"

	"github.com/maxatome/go-testdeep"
)

func TestMastParseNet(t *testing.T) {
	td := testdeep.NewT(t)
	res := mustParseNet("181.23.0.0/16")
	td.True(res.IP.Equal(net.IPv4(181, 23, 0, 0)))
	ones, size := res.Mask.Size()
	td.True(ones == 16 && size == 32)

	td.CmpPanic(func() {
		mustParseNet("asd")
	}, testdeep.NotNil())
}

func TestIsPublicIp(t *testing.T) {
	td := testdeep.NewT(t)
	td.True(isPublicIP(net.ParseIP("8.8.8.8")))
	td.True(isPublicIP(net.ParseIP("2a02:6b8:0:1::feed:0ff")))
	td.False(isPublicIP(net.ParseIP("")))
	td.False(isPublicIP(net.ParseIP("127.0.0.1")))
	td.False(isPublicIP(net.ParseIP("169.254.2.3")))
	td.False(isPublicIP(net.ParseIP("192.168.1.1")))
	td.False(isPublicIP(net.ParseIP("10.4.5.6")))
	td.False(isPublicIP(net.ParseIP("172.16.33.2")))
	td.False(isPublicIP(net.ParseIP("::")))
	td.False(isPublicIP(net.ParseIP("::1")))
	td.False(isPublicIP(net.ParseIP("::ffff:192.168.0.1")))
	td.False(isPublicIP(net.ParseIP("2001:db8::123")))
	td.False(isPublicIP(net.ParseIP("fe80::33")))
	td.False(isPublicIP(net.ParseIP("FC00::4")))
	td.False(isPublicIP(net.ParseIP("ff00::a")))
	td.False(isPublicIP(net.ParseIP("FF02:0:0:0:0:1:FF00::441")))
}

func TestGetBindedIpAddress(t *testing.T) {
	ctx, cancel := th.TestContext(t)
	defer cancel()

	td := testdeep.NewT(t)

	var f InterfacesAddrFunc = func() (addrs []net.Addr, e error) {
		return []net.Addr{
			&net.IPNet{IP: net.ParseIP("127.0.0.1"), Mask: net.CIDRMask(8, 32)},
			&net.IPNet{IP: net.ParseIP("161.32.6.19"), Mask: net.CIDRMask(32, 32)},
			&net.IPNet{IP: net.ParseIP("::1"), Mask: net.CIDRMask(128, 128)},
			&net.IPNet{IP: net.ParseIP("1.2.3.4"), Mask: net.CIDRMask(32, 32)},
			&net.IPNet{IP: net.ParseIP(""), Mask: net.CIDRMask(32, 32)},
			&net.IPNet{IP: net.ParseIP("2a02:6b8::feed:0ff"), Mask: net.CIDRMask(64, 128)},
		}, nil
	}

	res := getBindedIPAddress(ctx, f)
	td.CmpDeeply(res, []net.IP{net.ParseIP("127.0.0.1"), net.ParseIP("161.32.6.19"), net.ParseIP("::1"),
		net.ParseIP("1.2.3.4"),
		net.ParseIP("2a02:6b8::feed:0ff")})
}

func TestFilterPublicOnlyIPs(t *testing.T) {
	td := testdeep.NewT(t)

	ips := []net.IP{net.ParseIP("127.0.0.1"), net.ParseIP("161.32.6.19"), net.ParseIP("::1"),
		net.ParseIP("1.2.3.4"), net.ParseIP("2a02:6b8::feed:0ff")}

	res := filterPublicOnlyIPs(ips)
	td.CmpDeeply(res, []net.IP{net.ParseIP("161.32.6.19"), net.ParseIP("1.2.3.4"),
		net.ParseIP("2a02:6b8::feed:0ff")})
}

func TestIPList_Update(t *testing.T) {
	ctx, cancel := th.TestContext(t)
	defer cancel()

	td := testdeep.NewT(t)

	var retError = false
	l := NewIPList(ctx, func(ctx context.Context) (ips []net.IP, e error) {
		if retError {
			return nil, errors.New("test")
		}
		return nil, nil
	})

	l.ctx = zc.WithLogger(context.Background(), zap.NewNop())
	retError = false
	td.CmpNotPanic(func() {
		l.updateIPs()
	})

	retError = true
	td.CmpNotPanic(func() {
		l.updateIPs()
	})

	retError = false
	td.CmpNotPanic(func() {
		l.updateIPs()
	})

	l.ctx = zc.WithLogger(context.Background(), zap.NewNop().WithOptions(zap.Development()))
	retError = false
	td.CmpNotPanic(func() {
		l.updateIPs()
	})

	retError = true
	td.CmpPanic(func() {
		l.updateIPs()
	}, testdeep.NotNil())

	retError = false
	td.CmpNotPanic(func() {
		l.updateIPs()
	})
}

func TestIPList_UpdateByTimer(t *testing.T) {
	ctx, cancel := th.TestContext(t)
	defer cancel()

	td := testdeep.NewT(t)

	var ctxPanic, ctxPanicCancel = context.WithCancel(context.Background())

	var callCounter = 0

	var f AllowedIPAddresses = func(ctx context.Context) (addrs []net.IP, e error) {
		callCounter++
		if callCounter <= 2 {
			return []net.IP{
				net.ParseIP("1.2.3.4"),
			}, nil
		}

		if ctxPanic.Err() == nil {
			return []net.IP{
				net.ParseIP("161.32.6.19"), net.ParseIP("1.2.3.4"), net.ParseIP("2a02:6b8::feed:0ff"),
			}, nil
		}
		panic("must not called")
	}

	s := NewIPList(ctx, f)
	s.Addresses = f
	s.AutoUpdateInterval = 10 * time.Millisecond

	s.mu.RLock()
	td.CmpDeeply(len(s.ips), 1)
	td.True(s.ips[0].Equal(net.ParseIP("1.2.3.4")))
	s.mu.RUnlock()

	s.StartAutoRenew()

	time.Sleep(50 * time.Millisecond)

	s.mu.RLock()
	td.CmpDeeply(len(s.ips), 3)
	td.True(s.ips[0].Equal(net.ParseIP("161.32.6.19")))
	td.True(s.ips[1].Equal(net.ParseIP("1.2.3.4")))
	td.True(s.ips[2].Equal(net.ParseIP("2a02:6b8::feed:0ff")))
	s.mu.RUnlock()

	cancel()

	time.Sleep(50 * time.Millisecond)

	ctxPanicCancel()

	time.Sleep(50 * time.Millisecond)
}

func TestSetDefaultResolver(t *testing.T) {
	oldResolver := defaultResolver
	defer func() { // nolint:wsl
		defaultResolver = oldResolver
	}()

	resolver := NewResolverMock(t)
	SetDefaultResolver(resolver)
	testdeep.CmpDeeply(t, defaultResolver, resolver)
}

func TestSelfPublicIP_IsDomainAllowed(t *testing.T) {
	var _ DomainChecker = &IPList{}

	ctx, cancel := th.TestContext(t)
	defer cancel()

	ctx2, ctx2Cancel := th.TestContext(t)
	defer ctx2Cancel()

	td := testdeep.NewT(t)
	resolver := NewResolverMock(td)
	defer resolver.MinimockFinish()

	var res bool
	var err error

	s := NewIPList(ctx, func(ctx context.Context) (ips []net.IP, e error) {
		return []net.IP{net.ParseIP("1.2.3.4"), net.ParseIP("::ffff:2.2.2.2"), net.ParseIP("::1234")}, nil
	})

	s.Resolver = resolver
	resolver.LookupIPAddrMock.Expect(ctx2, "asd").Return([]net.IPAddr{{IP: net.ParseIP("1.2.3.4")}}, nil)
	res, err = s.IsDomainAllowed(ctx2, "asd")
	td.True(res)
	td.CmpNoError(err)

	resolver.LookupIPAddrMock.Expect(ctx2, "asd2").Return([]net.IPAddr{{IP: net.ParseIP("2.2.2.2")}}, nil)
	res, err = s.IsDomainAllowed(ctx2, "asd2")
	td.True(res)
	td.CmpNoError(err)

	resolver.LookupIPAddrMock.Expect(ctx2, "asd3").Return([]net.IPAddr{{IP: net.ParseIP("::1234")}}, nil)
	res, err = s.IsDomainAllowed(ctx2, "asd3")
	td.True(res)
	td.CmpNoError(err)

	resolver.LookupIPAddrMock.Expect(ctx2, "asd4").Return([]net.IPAddr{{IP: net.ParseIP("2.3.4.5")}}, nil)
	res, err = s.IsDomainAllowed(ctx2, "asd4")
	td.False(res)
	td.CmpNoError(err)

	resolver.LookupIPAddrMock.Expect(ctx2, "asd5").Return([]net.IPAddr{{IP: net.ParseIP("2.2.2.2")}},
		errors.New("test"))
	res, err = s.IsDomainAllowed(ctx2, "asd5")
	td.False(res)
	td.CmpError(err)

	resolver.LookupIPAddrMock.Expect(ctx2, "asd6").Return([]net.IPAddr{
		{IP: net.ParseIP("1.2.3.4")}, {IP: net.ParseIP("::1234")},
	}, nil)
	res, err = s.IsDomainAllowed(ctx2, "asd6")
	td.True(res)
	td.CmpNoError(err)

	resolver.LookupIPAddrMock.Expect(ctx2, "asd7").Return([]net.IPAddr{
		{IP: net.ParseIP("1.2.3.4")}, {IP: net.ParseIP("::1:1234")},
	}, nil)
	res, err = s.IsDomainAllowed(ctx2, "asd7")
	td.False(res)
	td.CmpNoError(err)

	resolver.LookupIPAddrMock.Expect(ctx2, "asd8").Return(nil, errors.New("test"))
	res, err = s.IsDomainAllowed(ctx2, "asd8")
	td.False(res)
	td.CmpError(err)
}

func TestSelfPublicIP_IsDomainAllowed_CanceledMainContext(t *testing.T) {
	ctx, cancel := th.TestContext(t)
	defer cancel()

	mainCtx, mainCtxCancel := context.WithCancel(context.Background())
	mainCtx = zc.WithLogger(mainCtx, zap.NewNop())
	mainCtxCancel()

	td := testdeep.NewT(t)

	s := NewIPList(mainCtx, func(ctx context.Context) (ips []net.IP, e error) {
		return nil, nil
	})
	res, err := s.IsDomainAllowed(ctx, "asd")
	td.False(res)
	td.CmpError(err)
}

func TestIPList_DoubleStart(t *testing.T) {
	ctx, cancel := th.TestContext(t)
	defer cancel()

	td := testdeep.NewT(t)

	td.CmpPanic(func() {
		s := NewIPList(ctx, func(ctx context.Context) (ips []net.IP, e error) {
			return nil, nil
		})
		s.ctx = zc.WithLogger(ctx, zap.NewNop().WithOptions(zap.Development())) // force panic on dpanic

		s.StartAutoRenew()
		s.StartAutoRenew()
	}, testdeep.NotNil())

	td.CmpNotPanic(func() {
		s := NewIPList(ctx, func(ctx context.Context) (ips []net.IP, e error) {
			return nil, nil
		})
		s.ctx = zc.WithLogger(ctx, zap.NewNop()) // force no panic on dpanic

		s.StartAutoRenew()
		s.StartAutoRenew()
	})
}

func TestTruncatedCopyIPs(t *testing.T) {
	td := testdeep.NewT(t)
	td.Nil(truncatedCopyIPs(nil))
	td.Nil(truncatedCopyIPs(make([]net.IP, 0, 10)))

	res := truncatedCopyIPs([]net.IP{nil})
	td.CmpDeeply(res, []net.IP{nil})
	td.CmpDeeply(cap(res), 1)

	res = truncatedCopyIPs(make([]net.IP, 2, 10))
	td.CmpDeeply(res, []net.IP{nil, nil})
	td.CmpDeeply(cap(res), 2)
}
