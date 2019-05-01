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
	td.True(isPublicIp(net.ParseIP("8.8.8.8")))
	td.True(isPublicIp(net.ParseIP("2a02:6b8:0:1::feed:0ff")))
	td.False(isPublicIp(net.ParseIP("")))
	td.False(isPublicIp(net.ParseIP("127.0.0.1")))
	td.False(isPublicIp(net.ParseIP("169.254.2.3")))
	td.False(isPublicIp(net.ParseIP("192.168.1.1")))
	td.False(isPublicIp(net.ParseIP("10.4.5.6")))
	td.False(isPublicIp(net.ParseIP("172.16.33.2")))
	td.False(isPublicIp(net.ParseIP("::")))
	td.False(isPublicIp(net.ParseIP("::1")))
	td.False(isPublicIp(net.ParseIP("::ffff:​192.168.0.1")))
	td.False(isPublicIp(net.ParseIP("2001:db8::123")))
	td.False(isPublicIp(net.ParseIP("fe80::33")))
	td.False(isPublicIp(net.ParseIP("FC00::4")))
	td.False(isPublicIp(net.ParseIP("ff00::a")))
	td.False(isPublicIp(net.ParseIP("FF02:0:0:0:0:1:FF00::441")))
}

func TestGetSelfPublicIPs(t *testing.T) {
	ctx, cancel := th.TestContext()
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

	res := getSelfPublicIPs(ctx, f)
	td.CmpDeeply(res, []net.IP{net.ParseIP("161.32.6.19"), net.ParseIP("1.2.3.4"),
		net.ParseIP("2a02:6b8::feed:0ff")})
}

func TestNewSelfIP_UpdateByTimer(t *testing.T) {
	ctx, cancel := th.TestContext()
	defer cancel()

	td := testdeep.NewT(t)

	var ctxPanic, ctxPanicCancel = context.WithCancel(context.Background())

	var firstTimeCalled = true

	var f InterfacesAddrFunc = func() (addrs []net.Addr, e error) {
		if firstTimeCalled {
			firstTimeCalled = false
			return []net.Addr{
				&net.IPNet{IP: net.ParseIP("127.0.0.1"), Mask: net.CIDRMask(8, 32)},
				&net.IPNet{IP: net.ParseIP("1.2.3.4"), Mask: net.CIDRMask(32, 32)},
			}, nil
		}

		if ctxPanic.Err() == nil {
			return []net.Addr{
				&net.IPNet{IP: net.ParseIP("127.0.0.1"), Mask: net.CIDRMask(8, 32)},
				&net.IPNet{IP: net.ParseIP("161.32.6.19"), Mask: net.CIDRMask(32, 32)},
				&net.IPNet{IP: net.ParseIP("::1"), Mask: net.CIDRMask(128, 128)},
				&net.IPNet{IP: net.ParseIP("1.2.3.4"), Mask: net.CIDRMask(32, 32)},
				&net.IPNet{IP: net.ParseIP(""), Mask: net.CIDRMask(32, 32)},
				&net.IPNet{IP: net.ParseIP("2a02:6b8::feed:0ff"), Mask: net.CIDRMask(64, 128)},
			}, nil
		}
		panic("must not called")
	}

	s := NewSelfIP(ctx)
	s.Addresses = f
	s.AutoUpdateInterval = 10 * time.Millisecond

	s.Start()
	s.mu.RLock()
	td.True(len(s.ips) == 1)
	td.True(s.ips[0].Equal(net.ParseIP("1.2.3.4")))
	s.mu.RUnlock()

	time.Sleep(50 * time.Millisecond)

	s.mu.RLock()
	td.True(len(s.ips) == 3)
	td.True(s.ips[0].Equal(net.ParseIP("161.32.6.19")))
	td.True(s.ips[1].Equal(net.ParseIP("1.2.3.4")))
	td.True(s.ips[2].Equal(net.ParseIP("2a02:6b8::feed:0ff")))
	s.mu.RUnlock()

	cancel()

	time.Sleep(50 * time.Millisecond)

	ctxPanicCancel()

	time.Sleep(50 * time.Millisecond)
}

func TestSelfPublicIP_IsDomainAllowed(t *testing.T) {
	var _ DomainChecker = &SelfPublicIP{}

	ctx, cancel := th.TestContext()
	defer cancel()

	ctx2, ctx2Cancel := th.TestContext()
	defer ctx2Cancel()

	td := testdeep.NewT(t)
	resolver := NewResolverMock(td)
	defer resolver.MinimockFinish()
	var res bool
	var err error

	s := NewSelfIP(ctx)
	s.Resolver = resolver
	s.ips = []net.IP{net.ParseIP("1.2.3.4"), net.ParseIP("::ffff:2.2.2.2"), net.ParseIP("::1234")}

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

}

func TestNewSelfIP_DoubleStart(t *testing.T) {
	ctx, cancel := th.TestContext()
	defer cancel()

	td := testdeep.NewT(t)

	s := NewSelfIP(ctx)
	s.Start()

	td.CmpPanic(func() {
		s.ctx = zc.WithLogger(ctx, zap.NewNop().WithOptions(zap.Development()))
		s.Start()
	}, testdeep.NotNil())

	td.CmpNotPanic(func() {
		s.ctx = ctx
		s.Start()
	})
}