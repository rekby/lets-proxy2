//nolint:golint
package domain_checker

import (
	"context"
	"errors"
	"net"
	"testing"

	"github.com/gojuno/minimock/v3"

	"github.com/maxatome/go-testdeep"
	"github.com/rekby/lets-proxy2/internal/th"
)

func TestConfig_CreateDomainCheckerEmpty(t *testing.T) {
	ctx, cancel := th.TestContext(t)
	defer cancel()

	td := testdeep.NewT(t)
	cfg := Config{}
	checker, err := cfg.CreateDomainChecker(ctx)
	td.CmpNoError(err)

	res, err := checker.IsDomainAllowed(ctx, "asd")
	td.True(res)
	td.CmpNoError(err)
}

func TestConfig_CreateDomainCheckerBadBlackList(t *testing.T) {
	ctx, cancel := th.TestContext(t)
	defer cancel()

	td := testdeep.NewT(t)
	cfg := Config{
		BlackList: "12(",
	}
	res, err := cfg.CreateDomainChecker(ctx)
	td.Nil(res)
	td.CmpError(err)
}

func TestConfig_CreateDomainCheckerBadWhiteList(t *testing.T) {
	ctx, cancel := th.TestContext(t)
	defer cancel()

	td := testdeep.NewT(t)
	cfg := Config{
		WhiteList: "12(",
	}
	res, err := cfg.CreateDomainChecker(ctx)
	td.Nil(res)
	td.CmpError(err)
}

func TestConfig_CreateDomainCheckerBlackListOnly(t *testing.T) {
	ctx, cancel := th.TestContext(t)
	defer cancel()

	td := testdeep.NewT(t)
	cfg := Config{
		BlackList: `.*\.com$`,
	}
	checker, err := cfg.CreateDomainChecker(ctx)
	td.CmpNoError(err)

	res, err := checker.IsDomainAllowed(ctx, "asd.com")
	td.False(res)
	td.CmpNoError(err)

	res, err = checker.IsDomainAllowed(ctx, "asd.ru")
	td.True(res)
	td.CmpNoError(err)
}

func TestConfig_CreateDomainCheckerWhiteListOnly(t *testing.T) {
	ctx, cancel := th.TestContext(t)
	defer cancel()

	td := testdeep.NewT(t)
	cfg := Config{
		WhiteList: `.*\.com$`,
	}
	checker, err := cfg.CreateDomainChecker(ctx)
	td.CmpNoError(err)

	res, err := checker.IsDomainAllowed(ctx, "asd.com")
	td.True(res)
	td.CmpNoError(err)

	res, err = checker.IsDomainAllowed(ctx, "asd.ru")
	td.True(res)
	td.CmpNoError(err)
}

func TestConfig_CreateDomainCheckerSelfIPOnly(t *testing.T) {
	ctx, cancel := th.TestContext(t)
	defer cancel()

	td := testdeep.NewT(t)
	mc := minimock.NewController(td)
	defer mc.Finish()

	resolver := NewResolverMock(mc)
	resolver.LookupIPAddrMock.When(ctx, "bad").Then([]net.IPAddr{{IP: net.ParseIP("127.0.0.1")}}, nil)
	resolver.LookupIPAddrMock.When(ctx, "ok").Then([]net.IPAddr{{IP: net.ParseIP("1.2.3.4")}}, nil)
	resolver.LookupIPAddrMock.Return(nil, errors.New("unknown domain"))

	cfg := Config{
		IPSelf:             true,
		IPSelfDetectMethod: "bind",
	}

	checker, err := cfg.CreateDomainChecker(ctx)
	td.CmpNoError(err)
	ipList := checker.(All)[1].(Any)[0].(Any)[0].(*IPList)

	ipList.mu.Lock()
	ipList.Resolver = resolver
	ipList.Addresses = func(ctx context.Context) (ips []net.IP, e error) {
		return []net.IP{net.ParseIP("1.2.3.4")}, nil
	}
	ipList.mu.Unlock()
	ipList.updateIPs()

	res, err := checker.IsDomainAllowed(ctx, "bad")
	td.False(res)
	td.CmpNoError(err)

	res, err = checker.IsDomainAllowed(ctx, "ok")
	td.True(res)
	td.CmpNoError(err)

	res, err = checker.IsDomainAllowed(ctx, "unknown")
	td.False(res)
	td.CmpError(err)
}

func TestConfig_CreateDomainCheckerWhitelist(t *testing.T) {
	ctx, cancel := th.TestContext(t)
	defer cancel()

	td := testdeep.NewT(t)
	mc := minimock.NewController(td)
	defer mc.Finish()

	resolver := NewResolverMock(mc)
	resolver.LookupIPAddrMock.When(ctx, "whitelist").Then([]net.IPAddr{{IP: net.ParseIP("3.3.3.3")}}, nil)
	resolver.LookupIPAddrMock.When(ctx, "unknown").Then(nil, errors.New("unknown domain"))

	cfg := Config{
		IPWhiteList: "2.3.4.5,3.3.3.3",
	}

	checker, err := cfg.CreateDomainChecker(ctx)
	td.CmpNoError(err)
	whiteIPList := checker.(All)[1].(Any)[0].(*IPList)

	whiteIPList.mu.Lock()
	whiteIPList.Resolver = resolver
	whiteIPList.mu.Unlock()
	whiteIPList.updateIPs()

	res, err := checker.IsDomainAllowed(ctx, "whitelist")
	td.True(res)
	td.CmpNoError(err)

	res, err = checker.IsDomainAllowed(ctx, "unknown")
	td.False(res)
	td.CmpError(err)
}

func TestConfig_CreateDomainCheckerComplex(t *testing.T) {
	ctx, cancel := th.TestContext(t)
	defer cancel()

	td := testdeep.NewT(t)
	mc := minimock.NewController(td)
	defer mc.Finish()

	resolver := NewResolverMock(mc)
	resolver.LookupIPAddrMock.When(ctx, "bad.test.com").Then([]net.IPAddr{{IP: net.ParseIP("127.0.0.1")}}, nil)
	resolver.LookupIPAddrMock.When(ctx, "ok.test.com").Then([]net.IPAddr{{IP: net.ParseIP("1.2.3.4")}}, nil)
	resolver.LookupIPAddrMock.When(ctx, "bad.ru").Then([]net.IPAddr{{IP: net.ParseIP("127.0.0.1")}}, nil)
	resolver.LookupIPAddrMock.When(ctx, "ok.ru").Then([]net.IPAddr{{IP: net.ParseIP("1.2.3.4")}}, nil)
	resolver.LookupIPAddrMock.When(ctx, "whitelist").Then([]net.IPAddr{{IP: net.ParseIP("3.3.3.3")}}, nil)
	resolver.LookupIPAddrMock.When(ctx, "unknown").Then(nil, errors.New("unknown domain"))

	cfg := Config{
		BlackList:          `.*\.com`,
		WhiteList:          `(.*\.)?test\.com`,
		IPSelf:             true,
		IPSelfDetectMethod: "bind",
		IPWhiteList:        "2.3.4.5,3.3.3.3",
	}

	checker, err := cfg.CreateDomainChecker(ctx)
	td.CmpNoError(err)

	selfIPList := checker.(All)[1].(Any)[0].(Any)[0].(*IPList)
	selfIPList.mu.Lock()
	selfIPList.Resolver = resolver
	selfIPList.Addresses = func(ctx context.Context) (ips []net.IP, e error) {
		return []net.IP{net.ParseIP("1.2.3.4")}, nil
	}
	selfIPList.mu.Unlock()
	selfIPList.updateIPs()

	whiteIPList := checker.(All)[1].(Any)[1].(*IPList)
	whiteIPList.mu.Lock()
	whiteIPList.Resolver = resolver
	whiteIPList.mu.Unlock()
	whiteIPList.updateIPs()

	res, err := checker.IsDomainAllowed(ctx, "any.com")
	td.False(res)
	td.CmpNoError(err)

	res, err = checker.IsDomainAllowed(ctx, "bad.test.com")
	td.False(res)
	td.CmpNoError(err)

	res, err = checker.IsDomainAllowed(ctx, "ok.test.com")
	td.True(res)
	td.CmpNoError(err)

	res, err = checker.IsDomainAllowed(ctx, "bad.ru")
	td.False(res)
	td.CmpNoError(err)

	res, err = checker.IsDomainAllowed(ctx, "ok.ru")
	td.True(res)
	td.CmpNoError(err)

	res, err = checker.IsDomainAllowed(ctx, "whitelist")
	td.True(res)
	td.CmpNoError(err)

	res, err = checker.IsDomainAllowed(ctx, "unknown")
	td.False(res)
	td.CmpError(err)
}
